package setting

import (
	"log/slog"

	"github.com/spf13/viper"
)

func SettingEnv() {
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.SetConfigFile(".env")
	// Lambda 환경에서는 .env 파일이 없으므로 실패해도 계속 진행 (AutomaticEnv로 OS 환경변수 사용)
	if err := viper.ReadInConfig(); err != nil {
		slog.Info("Config file not loaded, using environment variables", "error", err)
	}
}
