package elastic

import (
	"babblegraph/util/env"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/getsentry/sentry-go"
)

const (
	elasticsearchHostsKey    string = "ES_HOSTS"
	elasticsearchUsernameKey string = "ES_USERNAME"
	elasticsearchPasswordKey string = "ES_PASSWORD"

	elasticsearchMigrationHostsKey    string = "ES_MIGRATION_HOSTS"
	elasticsearchMigrationUsernameKey string = "ES_MIGRATION_USERNAME"
	elasticsearchMigrationPasswordKey string = "ES_MIGRATION_PASSWORD"
)

var (
	esClient        *elasticsearch.Client
	migrationClient *elasticsearch.Client

	currentEnvironmentName *env.Environment
)

func InitializeElasticsearchClientForEnvironment() error {
	if esClient != nil {
		panic("elasticsearch client is already initialized")
	}
	currentEnvironmentName = env.MustEnvironmentName().Ptr()
	cfg := elasticsearch.Config{
		Addresses: getAddressesForEnvironment(elasticsearchHostsKey),
		Username:  env.GetEnvironmentVariableOrDefault(elasticsearchUsernameKey, "elastic"),
		Password:  env.MustEnvironmentVariable(elasticsearchPasswordKey),
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				KeepAlive: 60 * time.Second,
				DualStack: true,
				Timeout:   15 * time.Second,
			}).DialContext,
			MaxIdleConns:          2,
			IdleConnTimeout:       0,
			TLSHandshakeTimeout:   0,
			MaxIdleConnsPerHost:   0,
			ExpectContinueTimeout: 5 * time.Second,
		},
	}
	var err error
	esClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return err
	}
	migrationConfig := elasticsearch.Config{
		Addresses: getAddressesForEnvironment(elasticsearchMigrationHostsKey),
		Username:  env.GetEnvironmentVariableOrDefault(elasticsearchMigrationUsernameKey, "elastic"),
		Password:  env.MustEnvironmentVariable(elasticsearchMigrationPasswordKey),
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				KeepAlive: 60 * time.Second,
				DualStack: true,
				Timeout:   15 * time.Second,
			}).DialContext,
			MaxIdleConns:          2,
			IdleConnTimeout:       0,
			TLSHandshakeTimeout:   0,
			MaxIdleConnsPerHost:   0,
			ExpectContinueTimeout: 5 * time.Second,
		},
	}
	migrationClient, err = elasticsearch.NewClient(migrationConfig)
	return err
}

func handleMigrationError(err error) {
	switch *currentEnvironmentName {
	case env.EnvironmentProd,
		env.EnvironmentStage:
		sentry.CaptureException(err)
	case env.EnvironmentLocal,
		env.EnvironmentLocalNoEmail,
		env.EnvironmentLocalTestEmail:
		log.Println(err.Error())
	}
}

func getAddressesForEnvironment(hostsKey string) []string {
	addressesUnsplit := env.MustEnvironmentVariable(hostsKey)
	return strings.Split(addressesUnsplit, ",")
}
