package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Tencent struct {
		SecretId  string `yaml:"secretId"`
		SecretKey string `yaml:"secretKey"`
	} `yaml:"tencent"`
	Domain struct {
		Domain    string `yaml:"domain"`
		SubDomain string `yaml:"subDomain"`
	} `yaml:"domain"`
	CheckInterval int `yaml:"checkInterval"`
	Email         struct {
		SMTPServer string `yaml:"smtpServer"`
		SMTPPort   int    `yaml:"smtpPort"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		Recipient  string `yaml:"recipient"`
	} `yaml:"email"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Validate() error {
	// Add validation logic if needed
	return nil
}
