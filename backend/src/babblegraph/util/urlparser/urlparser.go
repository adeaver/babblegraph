package urlparser

import "regexp"

var multipleSlashesRegex = regexp.MustCompile("/*")

type ParsedURL struct {
	Domain        string
	URLIdentifier string
}

func ParseURL(rawURL string) *ParsedURL {
	return &ParsedURL{}
}

type urlParts struct {
	Website string
	Page    string
	Params  string
}

func findURLParts(rawURL string) *urlParts {
	return nil
}
