package config

import (
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

var Config AppConfig

type AppConfig struct {
	Ignore        []string
	configDir     string
	hooksDir      string
	SkipSetGlobal bool `mapstructure:"skip_set_global"`
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not get user home dir: %s", err)
	}

	Config = AppConfig{}

	Config.configDir = path.Join(home, ".config", "ail")
	Config.hooksDir = path.Join(Config.configDir, "hooks")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(Config.configDir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return
		}
		log.Fatalf("could not read config: %s", err)
	}

	if err := viper.UnmarshalExact(&Config); err != nil {
		log.Fatalf("could not parse config: %s", err)
	}
}

func HookPath(pkg string) string {
	return path.Join(Config.hooksDir, pkg+".sh")
}

func IsIgnored(pkg string) bool {
	for _, p := range Config.Ignore {
		if pkg == p {
			return true
		}
	}
	return false
}
