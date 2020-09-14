package links

import (
	"fmt"
	"math"
	"net/url"
	"strings"
)

func GetLinkForURL(rawURL string) (*Link, error) {
	domain, cleanURL, err := getDomainAndCleanURL(rawURL)
	if err != nil {
		return nil, err
	}
	return &Link{
		Domain: *domain,
		URL:    *cleanURL,
	}, nil
}

func getDomainAndCleanURL(rawURL string) (*Domain, *string, error) {
	u, err := url.Parse(ensureScheme(rawURL))
	if err != nil {
		return nil, nil, err
	}
	d, err := getDomainForURL(u)
	if err != nil {
		return nil, nil, err
	}
	cleanURL := getCleanURLForURL(u)
	return d, &cleanURL, nil
}

func ensureScheme(rawURL string) string {
	switch {
	case strings.HasPrefix(rawURL, "http://"):
		return rawURL
	case strings.HasPrefix(rawURL, "https://"):
		return rawURL
	default:
		return fmt.Sprintf("http://%s", rawURL)
	}
}

func getDomainForURL(u *url.URL) (*Domain, error) {
	hostWithoutPort := u.Hostname()
	hostParts := strings.Split(hostWithoutPort, ".")
	if len(hostParts) == 0 {
		return nil, fmt.Errorf("invalid url")
	}
	endIndex := int(math.Max(float64(len(hostParts)-2), 0))
	filteredDomain := strings.Join(hostParts[endIndex:], ".")
	d := Domain(strings.ToLower(filteredDomain))
	return &d, nil
}

func getCleanURLForURL(u *url.URL) string {
	return fmt.Sprintf("%s://%s%s", u.Scheme, u.Hostname(), u.Path)
}
