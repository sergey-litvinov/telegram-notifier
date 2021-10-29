package main

import (
	"bytes"
	"errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"strings"
)

type Config struct{
	Telegram TelegramConfig
	Healthcheck HealthcheckConfig
}

type TelegramConfig struct{
	Token string
	ForwardTo int64
	Debug bool
}

type HealthcheckConfig struct{
	Hosts[] string
	Debug bool
}

//loadConfig loads configuration from the environment variables
func loadConfig() (Config, error){
	var config Config

	v := viper.New()
	// set default values in viper.
	// Viper needs to know if a key exists in order to override it.
	// https://github.com/spf13/viper/issues/188
	b, err := yaml.Marshal(&config)
	if err != nil {
		return Config{}, err
	}
	defaultConfig := bytes.NewReader(b)
	v.SetConfigType("yaml")
	if err := v.MergeConfig(defaultConfig); err != nil {
		return Config{}, err
	}
	// overwrite values from config
	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return Config{}, err
		}
		// don't return error if file is missing. overwrite file is optional
	}
	// tell viper to overwrite env variables
	v.AutomaticEnv()
	v.SetEnvPrefix("NOTIFIER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// refresh configuration with all merged values
	v.Unmarshal(&config)

	if len(config.Telegram.Token) == 0 {
		return Config{}, errors.New("telegram token env variable isn't found. please verify that NOTIFIER_TELEGRAM_TOKEN variable is present")
	}

	return config, nil
}

