package image

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
)

// govips를 사용해 EXIF 메타데이터 추출 (HEIC 포함)
func extractExif(imageBytes []byte) (takenAt time.Time, lat *float64, lon *float64) {
	takenAt = time.Now()

	img, err := vips.NewImageFromBuffer(imageBytes)
	if err != nil {
		return
	}
	defer img.Close()

	exifData := img.GetExif()

	// 촬영 일시: DateTimeOriginal 우선, 없으면 DateTime
	for _, key := range []string{"exif-ifd2-DateTimeOriginal", "exif-ifd0-DateTime"} {
		if raw, ok := exifData[key]; ok {
			dateStr := extractExifValue(raw)
			if t, err := time.Parse("2006:01:02 15:04:05", dateStr); err == nil {
				takenAt = t
				break
			}
		}
	}

	// GPS 좌표
	latStr, latOk := exifData["exif-ifd3-GPSLatitude"]
	lonStr, lonOk := exifData["exif-ifd3-GPSLongitude"]
	if !latOk || !lonOk {
		return
	}

	latVal, err := parseGPSRational(extractExifValue(latStr))
	if err != nil {
		return
	}
	lonVal, err := parseGPSRational(extractExifValue(lonStr))
	if err != nil {
		return
	}

	if ref, ok := exifData["exif-ifd3-GPSLatitudeRef"]; ok {
		if extractExifValue(ref) == "S" {
			latVal = -latVal
		}
	}
	if ref, ok := exifData["exif-ifd3-GPSLongitudeRef"]; ok {
		if extractExifValue(ref) == "W" {
			lonVal = -lonVal
		}
	}

	lat = &latVal
	lon = &lonVal
	return
}

// govips EXIF 값에서 실제 값 추출
// 예: "2025:11:26 13:33:59 (2025:11:26 13:33:59, ASCII, ...)" → "2025:11:26 13:33:59"
func extractExifValue(s string) string {
	s = strings.TrimRight(s, "\x00")
	if idx := strings.Index(s, " ("); idx >= 0 {
		return strings.TrimSpace(s[:idx])
	}
	return strings.TrimSpace(s)
}

// EXIF GPS rational 문자열 파싱
// 예: "37/1 28/1 456/100" → 십진수 도
func parseGPSRational(s string) (float64, error) {
	parts := strings.Fields(s)
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid GPS rational: %s", s)
	}

	divisors := []float64{1, 60, 3600}
	var result float64
	for i, part := range parts {
		sub := strings.SplitN(strings.TrimSpace(part), "/", 2)
		if len(sub) != 2 {
			return 0, fmt.Errorf("invalid rational: %s", part)
		}
		num, err1 := strconv.ParseFloat(sub[0], 64)
		den, err2 := strconv.ParseFloat(sub[1], 64)
		if err1 != nil || err2 != nil || den == 0 {
			return 0, fmt.Errorf("invalid rational components: %s", part)
		}
		result += (num / den) / divisors[i]
	}
	return result, nil
}
