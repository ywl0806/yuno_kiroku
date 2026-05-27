package utils

import (
	"fmt"
	"strings"
)

// BuildMediaKey returns "{prefix}/{familyId}/{mediaItemId}.{ext}"
// ext는 "." 없이 전달 (예: "jpg", "webp", "mp4")
func BuildMediaKey(familyId, prefix, mediaItemId, ext string) string {
	return fmt.Sprintf("%s/%s/%s.%s", prefix, familyId, mediaItemId, ext)
}

// BuildMediaKeyFromFileName fileName에서 확장자를 추출하여 BuildMediaKey를 호출
func BuildMediaKeyFromFileName(familyId, prefix, mediaItemId, fileName string) string {
	ext := strings.ToLower(strings.TrimPrefix(fileName[strings.LastIndex(fileName, "."):], "."))
	return BuildMediaKey(familyId, prefix, mediaItemId, ext)
}
