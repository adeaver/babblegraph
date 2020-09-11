package postgres

import (
	"fmt"

	"github.com/adeaver/babblegraph/lib/env"
)

type SSLModeOption string

const (
	SSLModeOptionRequire    SSLModeOption = "require"
	SSLModeOptionVerifyFull SSLModeOption = "verify-full"
	SSLModeOptionVerifyCA   SSLModeOption = "verify-ca"
	SSLModeOptionDisable    SSLModeOption = "disable"
)

func (o SSLModeOption) Str() string {
	return string(o)
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  SSLModeOption
}

func (c *PostgresConfig) MakeConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode.Str())
}

func MustPostgresConfigForEnvironment() PostgresConfig {
	return PostgresConfig{
		Host:     env.MustEnvironmentVariable("PG_HOST"),
		Port:     env.GetEnvironmentVariableOrDefault("PG_PORT", "5432"),
		User:     env.MustEnvironmentVariable("PG_USER"),
		Password: env.MustEnvironmentVariable("PG_PASSWORD"),
		DBName:   env.MustEnvironmentVariable("PG_DB_NAME"),
		SSLMode:  MustSSLModeOptionForString(env.GetEnvironmentVariableOrDefault("PG_SSL_MODE", SSLModeOptionDisable.Str())),
	}
}

func MustSSLModeOptionForString(o string) SSLModeOption {
	switch o {
	case "require":
		return SSLModeOptionRequire
	case "verify-full":
		return SSLModeOptionVerifyFull
	case "verify-ca":
		return SSLModeOptionVerifyCA
	case "disable":
		return SSLModeOptionDisable
	default:
		panic(fmt.Sprintf("unrecognized sslmode option: %s", o))
	}
}
