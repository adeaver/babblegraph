package urlparser

import (
	"babblegraph/util/deref"
	"babblegraph/util/ptr"
	"fmt"
	"regexp"
	"strings"
)

var (
	multipleSlashesRegex = regexp.MustCompile("/")
	pageSeparatorRegex   = regexp.MustCompile("\\?|#")
)

func EnsureProtocol(u string) (*string, error) {
	// This function is needed for Yahoo, which strips
	// hrefs if they don't include the protocol
	switch {
	case strings.HasPrefix(u, "http://"),
		strings.HasPrefix(u, "https://"):
		return ptr.String(u), nil
	case strings.HasPrefix(u, "ftp://"),
		strings.HasPrefix(u, "sftp://"),
		strings.HasPrefix(u, "mailto:"):
		return nil, fmt.Errorf("invalid prefix")
	default:
		return ptr.String(fmt.Sprintf("https://%s", u)), nil
	}
}

func IsValidURL(u string) bool {
	urlParts := findURLParts(u)
	if urlParts == nil {
		return false
	}
	verifiedDomain, _ := verifyDomain(urlParts.Website)
	if verifiedDomain == nil {
		return false
	}
	return true
}

type ParsedURL struct {
	Domain        string
	Path          string
	Protocol      *string
	URLIdentifier string
	URL           string
	Params        *string
}

func MustParseURL(rawURL string) ParsedURL {
	u := ParseURL(rawURL)
	if u == nil {
		panic("did not parse URL correctly")
	}
	return *u
}

func ParseURL(rawURL string) *ParsedURL {
	urlParts := findURLParts(rawURL)
	if urlParts == nil {
		return nil
	}
	verifiedDomain, verifiedSubdomain := verifyDomain(urlParts.Website)
	if verifiedDomain == nil {
		return nil
	}
	urlIdentifier := *verifiedDomain
	if len(urlParts.Page) > 0 {
		urlIdentifier = fmt.Sprintf("%s|%s", *verifiedDomain, urlParts.Page)
	}
	if len(*verifiedSubdomain) > 0 {
		urlIdentifier = fmt.Sprintf("%s|%s", *verifiedSubdomain, urlIdentifier)
	}
	var params *string
	if len(urlParts.Params) > 0 {
		params = ptr.String(urlParts.Params)
	}
	return &ParsedURL{
		Domain:        *verifiedDomain,
		Path:          urlParts.Page,
		Protocol:      urlParts.Protocol,
		URLIdentifier: urlIdentifier,
		URL:           rawURL,
		Params:        params,
	}
}

type urlParts struct {
	Website  string
	Page     string
	Params   string
	Protocol *string
}

func findURLParts(rawURL string) *urlParts {
	parts := multipleSlashesRegex.Split(rawURL, -1)
	var website, page, protocol *string
	var pageParts, paramParts []string
	for _, part := range parts {
		switch {
		case len(part) == 0:
			// This is a very common problem because
			// https:// causes there to be an empty string
			continue
		case website == nil && strings.Count(part, ".") > 0:
			// the first occurrence of a dot indicates a url
			// There is an edge case in which the URL parser doesn't work if there is a query string immediately after
			// the domain
			subParts := strings.Split(part, "?")
			if len(subParts) == 0 {
				// This shouldn't be possible, but I guess continue?
				continue
			}
			website = ptr.String(subParts[0])
			paramParts = append(paramParts, subParts[1:]...)
		case website != nil && page == nil:
			// now, we've captured the website name and are capturing pages
			if pageSeparatorRegex.MatchString(part) {
				pageSplit := pageSeparatorRegex.Split(part, 2)
				pageParts = append(pageParts, pageSplit[0])
				paramParts = append(paramParts, pageSplit[1])
				p := strings.Join(pageParts, "/")
				page = &p
				continue
			}
			pageParts = append(pageParts, part)
		case website != nil && page != nil:
			// now we're just collecting params
			paramParts = append(paramParts, part)
		default:
			// no-op
		}
	}
	if len(parts) > 0 && strings.HasPrefix(parts[0], "http") {
		protocolStr := parts[0]
		if !strings.HasSuffix(protocolStr, ":") {
			protocolStr = fmt.Sprintf("%s:", protocolStr)
		}
		protocol = ptr.String(protocolStr)
	}
	if website == nil {
		return nil
	}
	return &urlParts{
		Website:  *website,
		Page:     strings.TrimSuffix(deref.String(page, strings.Join(pageParts, "/")), "/"),
		Params:   strings.Join(paramParts, "&"),
		Protocol: protocol,
	}
}

func verifyDomain(website string) (_subdomain, _domain *string) {
	websiteParts := strings.Split(website, ".")
	var domainParts, tldParts []string
	for i := len(websiteParts) - 1; i >= 0; i-- {
		subParts := strings.Split(websiteParts[i], "?")
		if len(subParts) == 0 {
			continue
		}
		_, isValidTLD := validTLDs[subParts[0]]
		switch {
		case websiteParts[i] == "www":
			// We've gotten to what should be the top of the domain, exit.
			break
		case isValidTLD && len(domainParts) == 0:
			// We have yet to see a non-valid TLD part
			// so this must still be part of the TLD
			// NOTE: this algorithm has some nuance. It overindexes
			// on TLDs. So an ambiguous URL like "blog.musica.ar" should have
			// a url of "blog" and a tld of "musica.ar" - the idea here is that
			// there are, in theory, more websites with a multipart TLD than websites
			// with that domain (1)
			tldParts = append([]string{websiteParts[i]}, tldParts...) // prepend
		case isValidTLD && len(domainParts) > 0:
			// We've already seen a non-valid TLD part and a valid TLD part
			// so this is likely a subdomain that happens to be a valid TLD
			// like "mx.google.com" or something.
			domainParts = append([]string{websiteParts[i]}, domainParts...) // prepend
		case !isValidTLD && len(tldParts) > 0:
			// We've already seen some TLD parts, so this is probably the website name
			// or a subdomain
			domainParts = append([]string{websiteParts[i]}, domainParts...) // prepend
		case !isValidTLD && len(tldParts) == 0:
			// not a valid tld
			return nil, nil
		default:
			panic("unreachable")
		}
	}
	if len(tldParts) == 0 || len(tldParts) == 1 && len(domainParts) == 0 {
		// This is only possible if we see only a TLD
		// or www.TLD or www
		return nil, nil
	}
	var partsAssociatedWithDomain, partsAssociatedWithSubdomain []string
	if len(domainParts) > 0 {
		partsAssociatedWithDomain = append(partsAssociatedWithDomain, domainParts[len(domainParts)-1])
		partsAssociatedWithSubdomain = domainParts[:len(domainParts)-1]
	}
	partsAssociatedWithDomain = append(partsAssociatedWithDomain, tldParts...)
	verifiedDomain := strings.Join(partsAssociatedWithDomain, ".")
	verifiedSubdomain := strings.Join(partsAssociatedWithSubdomain, ".")
	return &verifiedDomain, &verifiedSubdomain
}
