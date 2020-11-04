package elastic

import (
	"babblegraph/util/env"
	"net"
	"net/http"
	"strings"
	"time"

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
	return err
}

func getAddressesForEnvironment() []string {
	addressesUnsplit := env.MustEnvironmentVariable(elasticsearchHostsKey)
	return strings.Split(addressesUnsplit, ",")
}
