package wordsmith

import (
	"log"
	"math/rand"
	"time"

	"github.com/adeaver/babblegraph/lib/env"
	"github.com/adeaver/babblegraph/lib/postgres"
)

var db *sqlx.DB

const (
	maxRetries int = 10

	minSleepFuzzMilliseconds         int = 1000
	maxSleepFuzzAdditionMilliseconds int = 800
)

func MustSetupWordsmithForEnvironment() error {
	if db != nil {
		panic("database already initialized")
	}
	config := getWordsmithPostgresConfigForEnvironment()
	fuzzCalc := rand.New(rand.NewSource(time.Now().UnixNano()))
	var err error
	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", config.MakeConnectionString())
		if err == nil {
			return nil
		}
		log.Println("Got error: %s. Retrying...", err.Error())
		fuzzMilliseconds := time.Duration(minSleepFuzzMilliseconds + fuzzCalc.Intn(maxSleepFuzzAdditionMilliseconds))
		time.Sleep(fuzzMilliseconds * time.Millisecond)
	}
	return err
}

func getWordsmithPostgresConfigForEnvironment() postgres.PostgresConfig {
	return PostgresConfig{
		Host:     env.MustEnvironmentVariable("WORDSMITH_HOST"),
		Port:     env.GetEnvironmentVariableOrDefault("WORDSMITH_PORT", "5432"),
		User:     env.MustEnvironmentVariable("WORDSMITH_USER"),
		Password: env.MustEnvironmentVariable("WORDSMITH_PASSWORD"),
		DBName:   env.MustEnvironmentVariable("WORDSMITH_DB_NAME"),
		SSLMode:  mustSSLModeOptionForString(env.GetEnvironmentVariableOrDefault("WORDSMITH_PG_SSL_MODE", SSLModeOptionDisable.Str())),
	}
}

func withTx(func(tx *sqlx.Tx) error) error {
	tx := db.MustBegin()
	if err := f(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Fatalf("query failed: %s, unable to rollback: %s", err.Error(), rbErr.Error())
		}
		return err
	}
	return tx.Commit()
}
