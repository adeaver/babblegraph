package newsletter

import "babblegraph/model/domains"

func CreateNewsletter(accessor userPreferencesAccessor) (*Newsletter, error) {

}

func getAllowableDomains(accessor userPreferencesAccessor) ([]string, error) {
	currentUserDomainCounts, err := accessor.getUserDomainCounts()
	if err != nil {
		return nil, err
	}
	domainCountByDomain := make(map[string]int64)
	for _, domainCount := range currentUserDomainCounts {
		domainCountByDomain[domainCount.Domain] = domainCount.Count
	}
	var out []string
	for _, d := range domains.GetDomains() {
		countForDomain, ok := domainCountByDomain[d]
		if ok {
			metadata, err := domains.GetDomainMetadata(d)
			if err != nil {
				return nil, err
			}
			if metadata.NumberOfMonthlyFreeArticles != nil && countForDomain >= *metadata.NumberOfMonthlyFreeArticles {
				continue
			}
		}
		out = append(out, d)
	}
	return out, nil
}
