package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Server   ServerConfig   `yaml:"server"`
	CORS     CORSConfig     `yaml:"cors"`
	Redis    RedisConfig    `yaml:"redis"`
	Logging  LoggingConfig  `yaml:"logging"`
}

type AppConfig struct {
	Name  string `yaml:"name"`
	Env   string `yaml:"env"`
	Port  int    `yaml:"port"`
	Debug bool   `yaml:"debug"`
}

type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Name            string        `yaml:"name"`
	SSLMode         string        `yaml:"ssl_mode"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expiry int64  `yaml:"expiry"`
}

type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Enabled  bool   `yaml:"enabled"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

// Load loads configuration from config.yaml and environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{}

	// Try to load from config.yaml
	if _, err := os.Stat("config.yaml"); err == nil {
		data, err := os.ReadFile("config.yaml")
		if err != nil {
			return nil, fmt.Errorf("failed to read config.yaml: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config.yaml: %w", err)
		}
	}

	// Override with environment variables
	cfg.overrideFromEnv()

	return cfg, nil
}

// overrideFromEnv overrides configuration with environment variables
func (c *Config) overrideFromEnv() {
	if env := os.Getenv("APP_NAME"); env != "" {
		c.App.Name = env
	}
	if env := os.Getenv("APP_ENV"); env != "" {
		c.App.Env = env
	}
	if env := os.Getenv("APP_PORT"); env != "" {
		fmt.Sscanf(env, "%d", &c.App.Port)
	}
	if env := os.Getenv("APP_DEBUG"); env != "" {
		c.App.Debug = env == "true"
	}

	if env := os.Getenv("DB_HOST"); env != "" {
		c.Database.Host = env
	}
	if env := os.Getenv("DB_PORT"); env != "" {
		fmt.Sscanf(env, "%d", &c.Database.Port)
	}
	if env := os.Getenv("DB_USER"); env != "" {
		c.Database.User = env
	}
	if env := os.Getenv("DB_PASSWORD"); env != "" {
		c.Database.Password = env
	}
	if env := os.Getenv("DB_NAME"); env != "" {
		c.Database.Name = env
	}
	if env := os.Getenv("DB_SSL_MODE"); env != "" {
		c.Database.SSLMode = env
	}

	if env := os.Getenv("JWT_SECRET"); env != "" {
		c.JWT.Secret = env
	}
	if env := os.Getenv("JWT_EXPIRY"); env != "" {
		fmt.Sscanf(env, "%d", &c.JWT.Expiry)
	}

	if env := os.Getenv("SERVER_HOST"); env != "" {
		c.Server.Host = env
	}
	if env := os.Getenv("SERVER_PORT"); env != "" {
		fmt.Sscanf(env, "%d", &c.Server.Port)
	}
}

// IsDevelopment returns true if the app is in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

// IsProduction returns true if the app is in production mode
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}
