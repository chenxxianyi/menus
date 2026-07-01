package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Upload   UploadConfig   `mapstructure:"upload"`
	AI       AIConfig       `mapstructure:"ai"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Debug   bool   `mapstructure:"debug"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type DatabaseConfig struct {
	Host                   string `mapstructure:"host"`
	Port                   int    `mapstructure:"port"`
	User                   string `mapstructure:"user"`
	Password               string `mapstructure:"password"`
	DBName                 string `mapstructure:"dbname"`
	Charset                string `mapstructure:"charset"`
	Loc                    string `mapstructure:"loc"`
	MaxIdleConns           int    `mapstructure:"max_idle_conns"`
	MaxOpenConns           int    `mapstructure:"max_open_conns"`
	ConnMaxLifetimeMinutes int    `mapstructure:"conn_max_lifetime_minutes"`
	AutoMigrate            bool   `mapstructure:"auto_migrate"`
}

func (d DatabaseConfig) DSN() string {
	loc := d.Loc
	if loc == "" {
		loc = "Local"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.Charset, loc)
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type JWTConfig struct {
	SecretKey   string `mapstructure:"secret_key"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

func (j JWTConfig) ExpireDuration() time.Duration {
	return time.Duration(j.ExpireHours) * time.Hour
}

type CORSConfig struct {
	Origins []string `mapstructure:"origins"`
}

type UploadConfig struct {
	Dir     string `mapstructure:"dir"`
	MaxSize int64  `mapstructure:"max_size"`
}

type AIConfig struct {
	BaseURL     string  `mapstructure:"base_url"`
	APIKey      string  `mapstructure:"api_key"`
	Model       string  `mapstructure:"model"`
	TimeoutSecs int     `mapstructure:"timeout_secs"`
	Temperature float64 `mapstructure:"temperature"`
}

func LoadConfig() (*Config, error) {
	_ = gotenv.Load(".env")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	applyDatabaseEnv(&cfg.Database)
	if v := os.Getenv("JWT_SECRET_KEY"); v != "" {
		cfg.JWT.SecretKey = v
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		cfg.Redis.Password = v
	}
	applyAIEnv(&cfg.AI)

	return &cfg, nil
}

func applyAIEnv(cfg *AIConfig) {
	if v := os.Getenv("AI_BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	if v := os.Getenv("AI_API_KEY"); v != "" {
		cfg.APIKey = v
	}
	if v := os.Getenv("AI_MODEL"); v != "" {
		cfg.Model = v
	}
	if v := os.Getenv("AI_TIMEOUT_SECS"); v != "" {
		if timeoutSecs, err := strconv.Atoi(v); err == nil {
			cfg.TimeoutSecs = timeoutSecs
		}
	}
	if v := os.Getenv("AI_TEMPERATURE"); v != "" {
		if temperature, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.Temperature = temperature
		}
	}
}

func applyDatabaseEnv(cfg *DatabaseConfig) {
	if v := os.Getenv("DATABASE_HOST"); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv("DATABASE_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Port = port
		}
	}
	if v := os.Getenv("DATABASE_NAME"); v != "" {
		cfg.DBName = v
	}
	if v := os.Getenv("DATABASE_USER"); v != "" {
		cfg.User = v
	}
	if v := os.Getenv("DATABASE_PASSWORD"); v != "" {
		cfg.Password = v
	}
	if v := os.Getenv("DATABASE_CHARSET"); v != "" {
		cfg.Charset = v
	}
	if v := os.Getenv("DATABASE_LOC"); v != "" {
		cfg.Loc = v
	}
	if v := os.Getenv("DATABASE_MAX_OPEN_CONNS"); v != "" {
		if conns, err := strconv.Atoi(v); err == nil {
			cfg.MaxOpenConns = conns
		}
	}
	if v := os.Getenv("DATABASE_MAX_IDLE_CONNS"); v != "" {
		if conns, err := strconv.Atoi(v); err == nil {
			cfg.MaxIdleConns = conns
		}
	}
	if v := os.Getenv("DATABASE_CONN_MAX_LIFETIME_MINUTES"); v != "" {
		if minutes, err := strconv.Atoi(v); err == nil {
			cfg.ConnMaxLifetimeMinutes = minutes
		}
	}
	if v := os.Getenv("DB_AUTO_MIGRATE"); v != "" {
		if autoMigrate, err := strconv.ParseBool(v); err == nil {
			cfg.AutoMigrate = autoMigrate
		}
	}
}
