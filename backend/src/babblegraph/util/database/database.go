package database

import (
	"babblegraph/util/postgres"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

const (
	maxRetries int = 10

	minSleepFuzzMilliseconds         int = 1000
	maxSleepFuzzAdditionMilliseconds int = 800
)

func GetDatabaseForEnvironmentRetrying() error {
	if db != nil {
		panic("database already initialized")
	}
	fuzzCalc := rand.New(rand.NewSource(time.Now().UnixNano()))
	config := postgres.MustPostgresConfigForEnvironment()
	var err error
	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", config.MakeConnectionString())
		if err == nil {
			return nil
		}
		log.Println(fmt.Sprintf("Got error: %s. Retrying...", err.Error()))
		fuzzMilliseconds := time.Duration(minSleepFuzzMilliseconds + fuzzCalc.Intn(maxSleepFuzzAdditionMilliseconds))
		time.Sleep(fuzzMilliseconds * time.Millisecond)
	}
	return err
}

func WithTx(f func(tx *sqlx.Tx) error) error {
	tx := db.MustBegin()
	if err := f(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Fatalf("query failed: %s, unable to rollback: %s", err.Error(), rbErr.Error())
		}
		return err
	}
	return tx.Commit()
}
