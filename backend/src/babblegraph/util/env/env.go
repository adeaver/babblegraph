package env

import (
	"fmt"
	"log"
	"os"
)

func MustEnvironmentVariable(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Println("Hello")
	panic(fmt.Sprintf("environment not properly configured. no value for key %s", key))
}

func GetEnvironmentVariableOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
