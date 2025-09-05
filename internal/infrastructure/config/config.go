package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	instance *Config
	once     sync.Once
)

// NewConfig returns a singleton Config instance loaded from config.yaml.
func NewConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")                         // File name (without extension)
		viper.SetConfigType("yaml")                           // File type
		viper.AddConfigPath("internal/infrastructure/config") // Path to look for the config file
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}
		var cfg Config
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Fatalf("Unable to decode into struct: %v", err)
		}
		instance = &cfg
	})
	return instance
}

// Config holds all configuration for the auth service.
type Config struct {
	Argon2id   Argon2idConfig   `yaml:"argon2id"`
	Postgres   PostgresConfig   `yaml:"postgres"`
	Logging    LoggingConfig    `yaml:"logging"`
	JWT        JWTConfig        `yaml:"jwt"`
	Redis      RedisConfig      `yaml:"redis"`
	CSRF       CSRFConfig       `yaml:"csrf"`
	RateLimit  RateLimitConfig  `yaml:"rate_limit"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
	Service    ServiceConfig    `yaml:"service"`
	CORS       CORSConfig       `yaml:"cors"`
	Security   SecurityConfig   `yaml:"security"`
}

type Argon2idConfig struct {
	Time    uint32 `yaml:"time"`
	Memory  uint32 `yaml:"memory"`
	Threads uint8  `yaml:"threads"`
	KeyLen  uint32 `yaml:"key_length"`
	SaltLen int    `yaml:"salt_length"`
}

type PostgresConfig struct {
	Driver       string `yaml:"driver"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	SSLMode      string `yaml:"sslmode"`
	MigrationDir string `yaml:"migration_dir"`
}

type LoggingConfig struct {
	Level       string `yaml:"level"`
	Encoding    string `yaml:"encoding"`
	Output      string `yaml:"output"`
	ErrorOutput string `yaml:"error_output"`
}

type JWTConfig struct {
	AccessTokenExpiry  string `yaml:"access_token_expiry"`
	RefreshTokenExpiry string `yaml:"refresh_token_expiry"`
	Issuer             string `yaml:"issuer"`
	Audience           string `yaml:"audience"`
	SigningMethod      string `yaml:"signing_method"`
	PrivateKeyPath     string `yaml:"private_key_path"`
	PublicKeyPath      string `yaml:"public_key_path"`
}

type RedisConfig struct {
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	DB                 int    `yaml:"db"`
	Password           string `yaml:"password"`
	AccessTokenPrefix  string `yaml:"access_token_prefix"`
	RefreshTokenPrefix string `yaml:"refresh_token_prefix"`
}

type CSRFConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Header      string `yaml:"header"`
	TokenLength int    `yaml:"token_length"`
	Expiry      string `yaml:"expiry"`
}

type RateLimitConfig struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerMinute int  `yaml:"requests_per_minute"`
	Burst             int  `yaml:"burst"`
}

type MonitoringConfig struct {
	Enabled            bool   `yaml:"enabled"`
	PrometheusEndpoint string `yaml:"prometheus_endpoint"`
}

type ServiceConfig struct {
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	Env      string `yaml:"env"`
	Debug    bool   `yaml:"debug"`
	Timezone string `yaml:"timezone"`
	URL      struct {
		App string `yaml:"app"`
	} `yaml:"url"`
}

type CORSConfig struct {
	Enabled          bool     `yaml:"enabled"`
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	ExposedHeaders   []string `yaml:"exposed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

type SecurityConfig struct {
	ContentTypeNosniff    bool   `yaml:"content_type_nosniff"`
	FrameOptions          string `yaml:"frame_options"`
	ContentSecurityPolicy string `yaml:"content_security_policy"`
	XSSProtection         string `yaml:"xss_protection"`
	ReferrerPolicy        string `yaml:"referrer_policy"`
	HSTSMaxAge            int    `yaml:"hsts_max_age"`
	HSTSIncludeSubdomains bool   `yaml:"hsts_include_subdomains"`
}

func (d *PostgresConfig) BuildDsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", d.Host, d.Port, d.User, d.DBName, d.Password, d.SSLMode)
}
