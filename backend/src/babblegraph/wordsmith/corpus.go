package wordsmith

import "fmt"

const (
	// UPC Wiki Corpus
	// https://www.cs.upc.edu/~nlp/wikicorpus/tagged.es.tgz
	CorpusIDUPCWikiCorpus CorpusID = "escrp1upc-wiki-corpus"
)

func getCurrentCorpusIDForLanguageCode(languageCode LanguageCode) CorpusID {
	switch languageCode {
	case LanguageCodeSpanish:
		return CorpusIDUPCWikiCorpus
	default:
		panic(fmt.Sprintf("Unrecognized language code: %s", languageCode.Str()))
	}
}
