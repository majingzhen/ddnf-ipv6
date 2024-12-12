package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Tencent struct {
		SecretId  string `yaml:"secret_id"`
		SecretKey string `yaml:"secret_key"`
	}
	Domain struct {
		Domain    string `yaml:"domain"`
		SubDomain string `yaml:"sub_domain"`
	}
	CheckInterval int         `yaml:"check_interval"`
	Email         EmailConfig `yaml:"email"`
}

type EmailConfig struct {
	SMTPServer string `yaml:"smtp_server"`
	SMTPPort   int    `yaml:"smtp_port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	ToEmail    string `yaml:"to_email"`
}

// LoadConfig 从文件加载配置
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate 验证配置是否有效
func (c *Config) Validate() error {
	// TODO: 添加配置验证逻辑
	return nil
}
