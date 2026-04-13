// Package config
package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

var Config *Configuration

type Server struct {
	Name               string        `mapstructure:"name" yaml:"name"`
	Port               string        `mapstructure:"port" yaml:"port"`
	Mode               string        `mapstructure:"mode" yaml:"mode"`
	ServiceName        string        `mapstructure:"service_name" yaml:"service_name"`
	ReadTimeout        time.Duration `mapstructure:"readTimeout" yaml:"readTimeout"`
	WriteTimeout       time.Duration `mapstructure:"writeTimeout" yaml:"writeTimeout"`
	MaxHeaderMegabytes int           `mapstructure:"maxHeaderMegabytes" yaml:"maxHeaderMegabytes"`
}

type PostgreSQL struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     string `mapstructure:"port" yaml:"port"`
	User     string `mapstructure:"user" yaml:"user"`
	Pass     string `mapstructure:"pass" yaml:"pass"`
	Database string `mapstructure:"database" yaml:"database"`
	SSLMode  string `mapstructure:"sslmode" yaml:"sslmode"`
}

type Configuration struct {
	PostgreSQL PostgreSQL `mapstructure:"postgresql" yaml:"postgresql"`
	Log        Logs       `mapstructure:"log" yaml:"log"`
	Server     Server     `mapstructure:"server" yaml:"server"`
	Trace      Trace      `mapstructure:"trace" yaml:"trace"`
}

type Logs struct {
	Elastic   string `mapstructure:"elastic" yaml:"elastic"`
	EsLogin   string `mapstructure:"es_login" yaml:"es_login"`
	EsPass    string `mapstructure:"es_pass" yaml:"es_pass"`
	EsIndex   string `mapstructure:"es_index" yaml:"es_index"`
	Path      string `mapstructure:"path" yaml:"path"`
	EventPath string `mapstructure:"event_path" yaml:"event_path"`
}

type Trace struct {
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`
}

func Setup(configPath string) {
	var configuration Configuration

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	Config = &configuration
}

func GetConfig() *Configuration {
	return Config
}
