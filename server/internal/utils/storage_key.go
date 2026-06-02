package utils

import (
	"fmt"
	"strings"
)

// BuildOriginalKey returns "original/{familyId}/{mediaItemId}.{ext}"
// S3 PUT 이벤트 → SQS → Lambda 파이프라인이 original/ 경로를 감지하므로 변경 불가
func BuildOriginalKey(familyId, mediaItemId, ext string) string {
	return fmt.Sprintf("original/%s/%s.%s", familyId, mediaItemId, ext)
}

// BuildOriginalKeyFromFileName fileName에서 확장자를 추출하여 BuildOriginalKey를 호출
func BuildOriginalKeyFromFileName(familyId, mediaItemId, fileName string) string {
	ext := strings.ToLower(strings.TrimPrefix(fileName[strings.LastIndex(fileName, "."):], "."))
	return BuildOriginalKey(familyId, mediaItemId, ext)
}

// BuildMediaKey returns "media/{familyId}/{prefix}/{mediaItemId}.{ext}"
// CloudFront signed cookie는 media/{familyId}/* 패턴으로 family별 접근을 제한
// ext는 "." 없이 전달 (예: "jpg", "webp", "mp4")
func BuildMediaKey(familyId, prefix, mediaItemId, ext string) string {
	return fmt.Sprintf("media/%s/%s/%s.%s", familyId, prefix, mediaItemId, ext)
}

// BuildMediaKeyFromFileName fileName에서 확장자를 추출하여 BuildMediaKey를 호출
func BuildMediaKeyFromFileName(familyId, prefix, mediaItemId, fileName string) string {
	ext := strings.ToLower(strings.TrimPrefix(fileName[strings.LastIndex(fileName, "."):], "."))
	return BuildMediaKey(familyId, prefix, mediaItemId, ext)
}

// BuildIdentityKey returns "media/{familyId}/identities/{identityId}/face_{mediaItemId}.webp"
func BuildIdentityKey(familyId string, identityId int32, mediaItemId string) string {
	return fmt.Sprintf("media/%s/identities/%d/face_%s.webp", familyId, identityId, mediaItemId)
}
