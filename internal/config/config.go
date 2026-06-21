package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Name    string
	Env     string
	Version string
	Port    string
}

type PostgresConfig struct {
	ConnectionString   string
	SchemaName         string
	MaxLifetime        time.Duration
	MaxIdleTime        time.Duration
	MaxIdleConnection  int
	MaxOpenConnection  int
	SlowQueryThreshold time.Duration
}

type CORSConfig struct {
	AllowedOrigin   []string
	AllowedMethod   []string
	AllowedHeader   []string
	AllowCredential bool
	MaxAge          int
}

func NewConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:    Get("APPLICATION_NAME", "app"),
			Env:     Get("APPLICATION_ENV", "development"),
			Version: Get("APPLICATION_VERSION", "1.0.0"),
			Port:    Get("APPLICATION_PORT", "3000"),
		},
		Postgres: PostgresConfig{
			ConnectionString:   MustGet[string]("DB_CONNECTION_STRING"),
			SchemaName:         Get("DB_SCHEMA_NAME", "public"),
			MaxLifetime:        Get("DB_MAX_LIFETIME", 300*time.Second),
			MaxIdleTime:        Get("DB_MAX_IDLE_TIME", 300*time.Second),
			MaxIdleConnection:  Get("DB_MAX_IDLE_CONNECTION", 10),
			MaxOpenConnection:  Get("DB_MAX_OPEN_CONNECTION", 100),
			SlowQueryThreshold: Get("DB_SLOW_QUERY_THRESHOLD", 200*time.Second),
		},
		CORS: CORSConfig{
			AllowedOrigin:   Get("CORS_ALLOWED_ORIGIN", []string{"http://localhost:3000", "http://localhost:8080", "http://localhost:3033", "http://localhost:5000"}),
			AllowedMethod:   Get("CORS_ALLOWED_METHOD", []string{"GET", "POST", "PUT", "PATCH"}),
			AllowedHeader:   Get("CORS_ALLOWED_HEADER", []string{"Origin", "Content-Type", "Accpet", "Authorization"}),
			AllowCredential: Get("CORS_ALLOW_CREDENTIAL", true),
			MaxAge:          Get("CORS_MAX_AGE", 86400),
		},
	}
}

func Get[T any](key string, fallback T) T {
	stringValue, isExists := os.LookupEnv(key)
	if !isExists {
		return fallback
	}

	var envValue T
	var err error

	switch p := any(&envValue).(type) {
	case *string:
		*p = stringValue
	case *int:
		*p, err = strconv.Atoi(stringValue)
	case *bool:
		*p, err = strconv.ParseBool(stringValue)
	case *float64:
		*p, err = strconv.ParseFloat(stringValue, 64)
	case *time.Duration:
		*p, err = time.ParseDuration(stringValue)
	case *[]string:
		if stringValue == "" {
			*p = []string{}
		} else {
			*p = strings.Split(stringValue, ",")
		}
	default:
		panic(fmt.Sprintf("config.Get: unsupported type for key %q", key))
	}

	if err != nil {
		slog.Error("config.Get: failed to parse env value", "key", key, "value", stringValue, "error", err)
		return fallback
	}

	return envValue
}

func MustGet[T any](key string) T {
	stringValue, isExists := os.LookupEnv(key)
	if !isExists || stringValue == "" {
		panic(fmt.Sprintf("config.MustGet: required env %q is not set", key))
	}

	var envValue T
	var err error

	switch p := any(&envValue).(type) {
	case *string:
		*p = stringValue
	case *int:
		*p, err = strconv.Atoi(stringValue)
	case *bool:
		*p, err = strconv.ParseBool(stringValue)
	case *float64:
		*p, err = strconv.ParseFloat(stringValue, 64)
	case *time.Duration:
		*p, err = time.ParseDuration(stringValue)
	case *[]string:
		if stringValue == "" {
			*p = []string{}
		} else {
			*p = strings.Split(stringValue, ",")
		}
	default:
		panic(fmt.Sprintf("config.MustGet: unsupported type for key %q", key))
	}

	if err != nil {
		panic(fmt.Sprintf("config.MustGet: failed to parse env %q = %q: %v", key, stringValue, err))
	}

	return envValue
}
