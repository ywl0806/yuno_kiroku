package utils

import (
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/pkg/utils"
)

func ParseStoragePath(storageKey string) string {
	return utils.ParsePath(viper.GetString("MEDIA_URL"), storageKey)
}
