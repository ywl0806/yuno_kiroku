package setting

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func SettingEnv() {
	mode := os.Getenv("APP_MODE")
	viper.SetConfigType("env")

	modePrefix := ""

	if mode != "" {
		modePrefix = "." + mode
	}

	envFile := ".env" + modePrefix
	viper.SetConfigFile(envFile)
	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
}
