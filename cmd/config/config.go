package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var config Cfg

type Cfg struct {
	App   AppConfig `mapstructure:"app"`
	MySQL MySQL     `mapstructure:"mysql"`
}

type AppConfig struct {
	Image     string `mapstructure:"image"`
	ImageTag  string `mapstructure:"imageTag"`
	Name      string `mapstructure:"name"`
	Port      string `mapstructure:"port"`
	JWTSecret string `mapstructure:"jwtSecret"`
}

type MySQL struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

func Load() error {
	v := viper.New()

	// Defaults
	v.SetDefault("app.image", "richisntreal-backend")
	v.SetDefault("app.imageTag", "latest")
	v.SetDefault("app.name", "richisntreal")
	v.SetDefault("app.port", "8080")
	v.SetDefault("app.jwtSecret", "changeme")

	v.SetDefault("mysql.host", "localhost")
	v.SetDefault("mysql.port", "3306")
	v.SetDefault("mysql.username", "root")
	v.SetDefault("mysql.password", "")
	v.SetDefault("mysql.database", "richisntreal")

	// Env‑vars
	v.SetEnvPrefix("RISNTREAL")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Config file (optional)
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	if err := v.Unmarshal(&config); err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}
	return nil
}

func Get() Cfg {
	return config
}
