package config

import (
	"github.com/spf13/viper"
)

var configPath string

type EmailConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Site     SiteConfig     `mapstructure:"site"`
	Email    EmailConfig    `mapstructure:"email"`
}

type ServerConfig struct {
	Port          int      `mapstructure:"port"`
	Mode          string   `mapstructure:"mode"`
	TrustedProxies []string `mapstructure:"trusted_proxies"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireHour int    `mapstructure:"expire_hour"`
}

type UploadConfig struct {
	Path     string   `mapstructure:"path"`
	MaxSize  int64    `mapstructure:"max_size"`
	AllowExt []string `mapstructure:"allow_ext"`
}

type SiteConfig struct {
	Name     string `mapstructure:"name"`
	Brief    string `mapstructure:"brief"`
	PageSize int    `mapstructure:"page_size"`
}

var Global *Config

func Load(path string) (*Config, error) {
	configPath = path
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	Global = &cfg
	return &cfg, nil
}

func Save() error {
	if configPath == "" {
		return nil
	}
	viper.Set("site.name", Global.Site.Name)
	viper.Set("site.brief", Global.Site.Brief)
	viper.Set("site.page_size", Global.Site.PageSize)
	return viper.WriteConfigAs(configPath)
}
