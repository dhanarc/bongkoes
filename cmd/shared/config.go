package shared

import (
	"github.com/djk-lgtm/bongkoes/config"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	DBLocation     = "ENV_BONGKOES_DB_LOCATION"
	AtlassianEmail = "ENV_BONGKOES_ATLASSIAN_EMAIL"
	AtlassianToken = "ENV_BONGKOES_ATLASSIAN_TOKEN"
	ConfluenceHost = "ENV_BONGKOES_CONFLUENCE_HOST"
)

func InitConfig() *config.Config {
	cfg := initDefaultConfig()
	if cfg != nil {
		return cfg
	}

	return readEnv()
}

func readEnv() *config.Config {
	return &config.Config{
		Bongkoes: config.BongkoesConfig{
			DBLocation:     os.Getenv(DBLocation),
			AtlassianEmail: os.Getenv(AtlassianEmail),
			AtlassianToken: os.Getenv(AtlassianToken),
			ConfluenceHost: os.Getenv(ConfluenceHost),
		},
	}
}

func initDefaultConfig() *config.Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil
	}
	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil
	}
	return &cfg
}
