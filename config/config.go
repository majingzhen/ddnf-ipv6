package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Tencent struct {
		SecretId  string
		SecretKey string
	}
	Domain struct {
		Domain string

		SubDomain string
	}
	CheckInterval int
	Email         Email
	Proxy         struct {
		EnableHTTP      bool
		HTTPListenAddr  string
		HTTPTargetAddr  string
		EnableHTTPS     bool
		HTTPSListenAddr string
		HTTPSTargetAddr string
		CertFile        string
		KeyFile         string
	}
}

type Email struct {
	SMTPServer string
	SMTPPort   int
	Username   string
	Password   string
	Recipient  string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return &cfg, nil
}
