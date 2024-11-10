package shared

import (
	"github.com/djk-lgtm/bongkoes/config"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func InitConfig() *config.Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("invalid read config")
		return nil
	}
	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal("invalid read config")
		return nil
	}
	return &cfg
}
