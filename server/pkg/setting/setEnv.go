package setting

import (
	"log"

	"github.com/spf13/viper"
)

func SettingEnv() {
	viper.SetConfigType("env")

	viper.SetConfigFile(".env")
	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
}
