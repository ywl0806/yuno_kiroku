package video

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/enums"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	"github.com/ywl0806/yuno_kiroku/internal/utils"
	"github.com/ywl0806/yuno_kiroku/pkg/storage"
)

const maxVideoFileSize = 4 * 1024 * 1024 * 1024 // 4GB

type VideoProcessingService struct {
	log            *slog.Logger
	mediaItemStore store.MediaItemStore
	storage        storage.StorageService
}

func NewVideoProcessingService(
	mediaItemStore store.MediaItemStore,
	storageService storage.StorageService,
) *VideoProcessingService {
	return &VideoProcessingService{
		log:            slog.Default().With("layer", "worker", "component", "video"),
		mediaItemStore: mediaItemStore,
		storage:        storageService,
	}
}

// ProcessVideo SQS job 파라미터를 받아 비디오 처리 전체 파이프라인을 실행합니다.
func (s *VideoProcessingService) ProcessVideo(ctx context.Context, params services.VideoJobParams) error {
	tmpDir, err := os.MkdirTemp("", "yuno-video-"+params.MediaItemID[:8]+"-*")
	if err != nil {
		s.log.ErrorContext(ctx, "temp dir 생성 실패", "media_item_id", params.MediaItemID, "error", err)
		return err
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			s.log.WarnContext(ctx, "temp dir cleanup 실패 (무시)", "error", err)
		}
	}()

	inputPath := filepath.Join(tmpDir, "input.mp4")
	thumbPath := filepath.Join(tmpDir, "thumb.webp")
	outputPath := filepath.Join(tmpDir, "output.mp4")

	// 1. processing 상태로 변경
	if _, err = s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           params.MediaItemID,
		UploadStatus: string(enums.UploadStatusProcessing),
	}); err != nil {
		s.log.ErrorContext(ctx, "status processing 업데이트 실패", "media_item_id", params.MediaItemID, "error", err)
		return err
	}

	// S3-06: 다운로드 전 파일 크기 검증 (디스크 소진 방지)
	fileSize, err := s.storage.GetFileSize(ctx, params.OriginalStorageKey)
	if err != nil {
		s.log.WarnContext(ctx, "파일 크기 조회 실패 (무시하고 계속)", "error", err)
	} else if fileSize > maxVideoFileSize {
		cause := fmt.Errorf("파일 크기 초과: %d bytes", fileSize)
		s.log.ErrorContext(ctx, "파일 크기 초과", "media_item_id", params.MediaItemID, "size", fileSize, "max", maxVideoFileSize)
		s.setFailed(ctx, params.MediaItemID, cause)
		return cause
	}

	// 2. S3 스트리밍 다운로드 — S3-02: 에러 타입 분류 로깅
	s.log.InfoContext(ctx, "다운로드 시작", "media_item_id", params.MediaItemID, "storage_key", params.OriginalStorageKey)
	if err = s.storage.DownloadToFile(ctx, params.OriginalStorageKey, inputPath); err != nil {
		if strings.Contains(err.Error(), "NoSuchKey") || strings.Contains(err.Error(), "AccessDenied") {
			s.log.ErrorContext(ctx, "영구 실패 (IAM/파일미존재)", "media_item_id", params.MediaItemID, "error", err)
		} else {
			s.log.WarnContext(ctx, "일시적 실패 (네트워크)", "media_item_id", params.MediaItemID, "error", err)
		}
		s.setFailed(ctx, params.MediaItemID, err)
		return err
	}

	// 3. 비디오 리사이즈 (최대 1280px, H.264, faststart)
	s.log.InfoContext(ctx, "비디오 리사이즈 중", "media_item_id", params.MediaItemID)
	if err = resizeVideo(ctx, inputPath, outputPath); err != nil {
		s.setFailed(ctx, params.MediaItemID, err)
		return err
	}
	videoW, videoH, err := probeSize(outputPath)
	if err != nil {
		s.log.WarnContext(ctx, "비디오 치수 조회 실패 (무시)", "media_item_id", params.MediaItemID, "error", err)
	}

	// 4. 썸네일 추출 (리사이즈된 영상의 1초 지점, WebP)
	s.log.InfoContext(ctx, "썸네일 추출 중", "media_item_id", params.MediaItemID)
	if err = extractThumbnail(ctx, outputPath, thumbPath); err != nil {
		s.setFailed(ctx, params.MediaItemID, err)
		return err
	}
	thumbW, thumbH := videoW, videoH

	thumbKey := utils.BuildMediaKey(params.FamilyID, consts.THUMBNAIL_STORAGE_PREFIX, params.MediaItemID, "webp")
	videoKey := utils.BuildMediaKey(params.FamilyID, consts.VIDEO_STORAGE_PREFIX, params.MediaItemID, "mp4")

	// S3-04: DB INSERT 먼저 → S3 업로드 순서로 변경하여 고아 파일 방지
	// 5. media_files DB 저장 (thumbnail + video) — upsert로 SQS 재시도 시 멱등성 보장
	if _, err = s.mediaItemStore.UpsertMediaFile(ctx, db.UpsertMediaFileParams{
		MediaItemID: params.MediaItemID,
		Role:        string(enums.MediaItemRoleThumbnail),
		StorageKey:  thumbKey,
		Width:       sql.NullInt32{Int32: thumbW, Valid: thumbW > 0},
		Height:      sql.NullInt32{Int32: thumbH, Valid: thumbH > 0},
	}); err != nil {
		s.setFailed(ctx, params.MediaItemID, err)
		return err
	}
	if _, err = s.mediaItemStore.UpsertMediaFile(ctx, db.UpsertMediaFileParams{
		MediaItemID: params.MediaItemID,
		Role:        string(enums.MediaItemRoleVideo),
		StorageKey:  videoKey,
		Width:       sql.NullInt32{Int32: videoW, Valid: videoW > 0},
		Height:      sql.NullInt32{Int32: videoH, Valid: videoH > 0},
	}); err != nil {
		// thumbnail DB 레코드 롤백
		if rbErr := s.mediaItemStore.DeleteMediaFileByItemAndRole(ctx, params.MediaItemID, string(enums.MediaItemRoleThumbnail)); rbErr != nil {
			s.log.ErrorContext(ctx, "thumbnail DB 롤백 실패", "media_item_id", params.MediaItemID, "error", rbErr)
		}
		s.setFailed(ctx, params.MediaItemID, err)
		return err
	}

	// 6. 썸네일 S3 업로드
	if err = s.storage.UploadFromFile(ctx, thumbKey, "image/webp", thumbPath); err != nil {
		// DB 레코드 롤백
		s.rollbackMediaFiles(ctx, params.MediaItemID)
		s.setFailed(ctx, params.MediaItemID, err)
		return err
	}

	// 7. 처리된 비디오 S3 업로드
	if err = s.storage.UploadFromFile(ctx, videoKey, "video/mp4", outputPath); err != nil {
		// DB 레코드 롤백 + thumbnail S3 삭제
		s.rollbackMediaFiles(ctx, params.MediaItemID)
		if delErr := s.storage.DeleteFile(ctx, thumbKey); delErr != nil {
			s.log.ErrorContext(ctx, "thumbnail S3 롤백 실패", "media_item_id", params.MediaItemID, "error", delErr)
		}
		s.setFailed(ctx, params.MediaItemID, err)
		return err
	}

	// 8. upload_status = completed
	if _, err = s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           params.MediaItemID,
		UploadStatus: string(enums.UploadStatusCompleted),
	}); err != nil {
		s.log.WarnContext(ctx, "status completed 업데이트 실패 (무시)", "media_item_id", params.MediaItemID, "error", err)
	}

	s.log.InfoContext(ctx, "비디오 처리 완료", "media_item_id", params.MediaItemID)
	return nil
}

func (s *VideoProcessingService) setFailed(ctx context.Context, mediaItemID string, cause error) {
	if _, err := s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           mediaItemID,
		UploadStatus: string(enums.UploadStatusFailed),
	}); err != nil {
		s.log.ErrorContext(ctx, "status failed 업데이트 실패", "media_item_id", mediaItemID, "error", err)
	}
	s.log.ErrorContext(ctx, "비디오 처리 실패", "media_item_id", mediaItemID, "cause", cause)
}

func (s *VideoProcessingService) rollbackMediaFiles(ctx context.Context, mediaItemID string) {
	for _, role := range []string{string(enums.MediaItemRoleThumbnail), string(enums.MediaItemRoleVideo)} {
		if err := s.mediaItemStore.DeleteMediaFileByItemAndRole(ctx, mediaItemID, role); err != nil {
			s.log.ErrorContext(ctx, "media_file DB 롤백 실패", "media_item_id", mediaItemID, "role", role, "error", err)
		}
	}
}

// extractThumbnail 1초 지점에서 WebP 썸네일 1장을 추출합니다.
// exec.CommandContext 사용으로 context 취소 시 ffmpeg 프로세스 SIGKILL 보장 (S3-03)
func extractThumbnail(ctx context.Context, inputPath, outputPath string) error {
	tctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(tctx, "ffmpeg",
		"-ss", "1",
		"-i", inputPath,
		"-vframes", "1",
		"-vcodec", "libwebp",
		"-q:v", "75",
		"-y", outputPath,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("썸네일 추출 실패: %w, output: %s", err, string(out))
	}
	return nil
}

// probeSize ffprobe로 미디어 파일의 첫 번째 비디오 스트림 치수를 반환합니다.
func probeSize(path string) (width, height int32, err error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_streams",
		path,
	)
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("ffprobe 실패: %w", err)
	}
	var result struct {
		Streams []struct {
			Width  int32 `json:"width"`
			Height int32 `json:"height"`
		} `json:"streams"`
	}
	if err = json.Unmarshal(out, &result); err != nil {
		return 0, 0, fmt.Errorf("ffprobe JSON 파싱 실패: %w", err)
	}
	for _, s := range result.Streams {
		if s.Width > 0 && s.Height > 0 {
			return s.Width, s.Height, nil
		}
	}
	return 0, 0, fmt.Errorf("스트림에서 치수를 찾을 수 없음")
}

// resizeVideo 가로 최대 1280px(짝수 맞춤), H.264 CRF23, faststart로 재인코딩합니다.
// exec.CommandContext 사용으로 context 취소 시 ffmpeg 프로세스 SIGKILL 보장 (S3-03)
func resizeVideo(ctx context.Context, inputPath, outputPath string) error {
	tctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(tctx, "ffmpeg",
		"-i", inputPath,
		"-vf", "scale='min(1280,iw)':'-2',format=yuv420p",
		"-vcodec", "libx264",
		"-crf", "23",
		"-preset", "ultrafast",
		"-acodec", "aac",
		"-movflags", "+faststart",
		"-y", outputPath,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg 실패: %w, output: %s", err, string(out))
	}
	return nil
}
