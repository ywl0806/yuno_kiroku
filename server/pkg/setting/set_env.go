package setting

import (
	"log"

	"github.com/spf13/viper"
)

func SettingEnv() {
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.SetConfigFile(".env")
	// Find and read the config file
	// Lambda 환경에서는 .env 파일이 없으므로 실패해도 계속 진행 (AutomaticEnv로 OS 환경변수 사용)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not loaded (%s), using environment variables", err)
	}
}
