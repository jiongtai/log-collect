package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Kafka       `yaml:"kafka"`
	Es          `yaml:"es"`
	LogFilePath string `yaml:"logFile"`
	Etcd        `yaml:"etcd"`
}

type Kafka struct {
	Address     []string `yaml:"address"`
	ChanMaxSize int      `yaml:"chanMaxSize"`
}

type Es struct {
	Address string `yaml:"address"`
	Size    int    `yaml:"size"`
	Worker  int    `yaml:"worker"`
}

type Etcd struct {
	Address       string        `yaml:"address"`
	Timeout       time.Duration `yaml:"timeout"`
	LogCollectKey string        `yaml:"logCollectKey"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	var configViperConfig = viper.New()
	configViperConfig.SetConfigName("config")
	configViperConfig.SetConfigType("yaml")
	// todo：路径问题待解决，能否换成绝对路径？build的时候当前路径会变成cmd/app/agent。。
	configViperConfig.AddConfigPath("./config")
	if err := configViperConfig.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := configViperConfig.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
