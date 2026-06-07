package config

import (
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	AI       AIConfig       `mapstructure:"ai"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"readTimeout"`
	WriteTimeout int    `mapstructure:"writeTimeout"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxIdleConns int    `mapstructure:"maxIdleConns"`
	MaxOpenConns int    `mapstructure:"maxOpenConns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireHour int    `mapstructure:"expireHour"`
}

type AIConfig struct {
	DefaultProvider string           `mapstructure:"defaultProvider"`
	Providers       []ProviderConfig `mapstructure:"providers"`
}

type ProviderConfig struct {
	Name    string `mapstructure:"name"`
	APIKey  string `mapstructure:"apiKey"`
	BaseURL string `mapstructure:"baseUrl"`
	Model   string `mapstructure:"model"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge"`
}

var (
	cfg  *Config
	once sync.Once
)

func Init(configPath string) error {
	var err error
	once.Do(func() {
		viper.SetConfigFile(configPath)
		viper.AutomaticEnv()

		if err = viper.ReadInConfig(); err != nil {
			return
		}

		cfg = &Config{}
		err = viper.Unmarshal(cfg)
	})
	return err
}

func Get() *Config {
	return cfg
}
