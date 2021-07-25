package config

import (
	"github.com/jinzhu/configor"
)

var Config = struct {
	App struct {
		Env      string
		HTTPAddr string
		HTTPPort string
	}

	DB struct {
		Driver   string
		Host     string `default:"localhost"`
		Port     string `default:"3306"`
		Name     string
		User     string `default:"root"`
		Password string `required:"true"`
		Locale   string `default:"Asia/Jakarta"`
	}
}{}

func init() {
	configor.Load(&Config, "config.yaml")
}
