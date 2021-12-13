package env

import (
	"fmt"
	"os"
)

type Environment string

const (
	EnvironmentProd           Environment = "prod"
	EnvironmentStage          Environment = "stage"
	EnvironmentLocal          Environment = "local"
	EnvironmentLocalNoEmail   Environment = "local-no-emails"
	EnvironmentLocalTestEmail Environment = "local-test-emails"
	EnvironmentTest           Environment = "test"
)

func (e Environment) Str() string {
	return string(e)
}

func (e Environment) Ptr() *Environment {
	return &e
}

func mustEnvironmentFromString(s string) Environment {
	switch s {
	case "prod":
		return EnvironmentProd
	case "stage":
		return EnvironmentStage
	case "local":
		return EnvironmentLocal
	case "local-no-emails":
		return EnvironmentLocalNoEmail
	case "local-test-emails":
		return EnvironmentLocalTestEmail
	case "test":
		return EnvironmentTest
	default:
		panic(fmt.Sprintf("Unrecognized environment: %s", s))
	}
}

func MustEnvironmentVariable(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic(fmt.Sprintf("environment not properly configured. no value for key %s", key))
}

func GetEnvironmentVariableOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func MustEnvironmentName() Environment {
	return mustEnvironmentFromString(MustEnvironmentVariable("ENV"))
}
