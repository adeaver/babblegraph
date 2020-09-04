package database

import (
	"fmt"

	"babblegraph/worker/util"
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

type postgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  SSLModeOption
}

func (c *postgresConfig) makeConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode.Str())
}

func mustPostgresConfigForEnvironment() postgresConfig {
	return postgresConfig{
		Host:     util.MustEnvironmentVariable("PG_HOST"),
		Port:     util.GetEnvironmentVariableOrDefault("PG_PORT", "5432"),
		User:     util.MustEnvironmentVariable("PG_USER"),
		Password: util.MustEnvironmentVariable("PG_PASSWORD"),
		DBName:   util.MustEnvironmentVariable("PG_DB_NAME"),
		SSLMode:  mustSSLModeOptionForString(util.GetEnvironmentVariableOrDefault("PG_SSL_MODE", SSLModeOptionDisable.Str())),
	}
}

func mustSSLModeOptionForString(o string) SSLModeOption {
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
