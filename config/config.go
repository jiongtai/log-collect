package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Kafka `yaml:"kafka"`
	Es    `yaml:"es"`
}

type Kafka struct {
	Address     string `yaml:"address"`
	ChanMaxSize int    `yaml:"chanMaxSize"`
}

type Es struct {
	Address string `yaml:"address"`
	Size    int    `yaml:"size"`
	Worker  int    `yaml:"worker"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	var configViperConfig = viper.New()
	configViperConfig.SetConfigName("config")
	configViperConfig.SetConfigType("yaml")
	configViperConfig.AddConfigPath("./config")
	if err := configViperConfig.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := configViperConfig.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
