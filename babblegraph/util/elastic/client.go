package elastic

import (
	"babblegraph/util/env"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

const (
	elasticsearchHostsKey    string = "ES_HOSTS"
	elasticsearchUsernameKey string = "ES_USERNAME"
	elasticsearchPasswordKey string = "ES_PASSWORD"
)

var esClient *elasticsearch.Client

func InitializeElasticsearchClientForEnvironment() error {
	if esClient != nil {
		panic("elasticsearch client is already initialized")
	}
	cfg := elasticsearch.Config{
		Addresses: getAddressesForEnvironment(),
		Username:  env.GetEnvironmentVariableOrDefault(elasticsearchUsernameKey, "elastic"),
		Password:  env.MustEnvironmentVariable(elasticsearchPasswordKey),
	}
	var err error
	esClient, err = elasticsearch.NewClient(cfg)
	return err
}

func getAddressesForEnvironment() []string {
	addressesUnsplit := env.MustEnvironmentVariable(elasticsearchHostsKey)
	return strings.Split(addressesUnsplit, " ")
}
