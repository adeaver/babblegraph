package urlparser

import (
	"babblegraph/util/deref"
	"fmt"
	"regexp"
	"strings"
)

var multipleSlashesRegex = regexp.MustCompile("/")

type ParsedURL struct {
	Domain        string
	URLIdentifier string
}

func ParseURL(rawURL string) *ParsedURL {
	urlParts := findURLParts(rawURL)
	if urlParts == nil {
		return nil
	}
	verifiedDomain := verifyDomain(urlParts.Website)
	if verifiedDomain == nil {
		return nil
	}
	urlIdentifier := *verifiedDomain
	if len(urlParts.Page) > 0 {
		urlIdentifier = fmt.Sprintf("%s|%s", *verifiedDomain, urlParts.Page)
	}
	return &ParsedURL{
		Domain:        *verifiedDomain,
		URLIdentifier: urlIdentifier,
	}
}

type urlParts struct {
	Website string
	Page    string
	Params  string
}

func findURLParts(rawURL string) *urlParts {
	parts := multipleSlashesRegex.Split(rawURL, -1)
	var website, page *string
	var pageParts, paramParts []string
	for _, part := range parts {
		switch {
		case len(part) == 0:
			// This is a very common problem because
			// https:// causes there to be an empty string
			continue
		case website == nil && strings.Count(part, ".") > 0:
			// the first occurrence of a dot indicates a url
			p := part
			website = &p
		case website != nil && page == nil:
			// now, we've captured the website name and are capturing pages
			if strings.Count(part, "?") != 0 {
				pageSplit := strings.Split(part, "?")
				pageParts = append(pageParts, pageSplit[0])
				paramParts = append(paramParts, pageSplit[1:]...)
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
	if website == nil {
		return nil
	}
	return &urlParts{
		Website: *website,
		Page:    deref.String(page, strings.Join(pageParts, "/")),
		Params:  strings.Join(paramParts, "&"),
	}
}

func verifyDomain(website string) *string {
	websiteParts := strings.Split(website, ".")
	var validParts []string
	for i := len(websiteParts) - 1; i >= 0; i-- {
		_, isValidTLD := validTLDs[websiteParts[i]]
		switch {
		case websiteParts[i] == "www":
			// We've gotten to what should be the top of the domain, exit.
			break
		case isValidTLD,
			!isValidTLD && len(validParts) > 0:
			// we either have a valid TLD or we've already seen one,
			// we add it to the domain
			validParts = append([]string{websiteParts[i]}, validParts...)
		case !isValidTLD && len(validParts) == 0:
			// not a valid tld
			return nil
		default:
			panic("unreachable")
		}
	}
	if len(validParts) <= 1 {
		return nil
	}
	verifiedDomain := strings.Join(validParts, ".")
	return &verifiedDomain
}
