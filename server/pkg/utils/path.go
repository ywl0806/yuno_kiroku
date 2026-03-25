package utils

import (
	"strings"
)

func ParsePath(baseUrl string, paths ...string) string {
	url := baseUrl
	for _, path := range paths {
		isPathStartWithSlash := strings.HasPrefix(path, "/")
		isUrlEndWithSlash := strings.HasSuffix(url, "/")

		if isPathStartWithSlash && isUrlEndWithSlash {
			url += path[1:]
		} else if isPathStartWithSlash != isUrlEndWithSlash {
			url += path
		} else {
			url += "/" + path
		}
	}
	return url
}
