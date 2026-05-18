package consts

const (
	VIEW_MAX_LENGTH      = 2048
	THUMBNAIL_MAX_LENGTH = 512

	// MAX_IMAGE_FILE_SIZE 이미지 파일 크기 상한 (50 MB) — 초과 시 OOM 위험
	MAX_IMAGE_FILE_SIZE int64 = 50 * 1024 * 1024

	VIEW_STORAGE_PREFIX      = "view"
	THUMBNAIL_STORAGE_PREFIX = "thumbnail"
	ORIGINAL_STORAGE_PREFIX  = "original"
	VIDEO_STORAGE_PREFIX     = "video"

	IDENTITY_STORAGE_PREFIX = "identities"

	// 얼굴 검색 임베딩 유사도 임계값
	FACE_SEARCH_THRESHOLD = 0.6
)
