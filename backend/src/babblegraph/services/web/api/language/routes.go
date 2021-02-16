package language

import (
	"babblegraph/services/web/router"
	"babblegraph/wordsmith"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

func RegisterRouteGroups() error {
	return router.RegisterRouteGroup(router.RouteGroup{
		Prefix: "language",
		Routes: []router.Route{
			{
				Path:    "get_lemmas_matching_text_1",
				Handler: handleGetLemmasMatchingText,
			},
		},
	})
}

type getLemmasMatchingTextRequest struct {
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
	Token        string                 `json:"token"`
	Text         string                 `json:"text"`
}

type getLemmasMatchingTextResponse struct {
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
	Text         string                 `json:"text"`
	Lemmas       []lemma                `json:"lemmas"`
}

type lemma struct {
	Text         string            `json:"text"`
	ID           wordsmith.LemmaID `json:"id"`
	PartOfSpeech partOfSpeech      `json:"part_of_speech"`
	Defintions   []definition      `json:"definitions"`
}

type partOfSpeech struct {
	ID   wordsmith.PartOfSpeechID `json:"id"`
	Name string                   `json:"name"`
}

type definition struct {
	Text      string  `json:"text"`
	ExtraInfo *string `json:"extra_info,omitempty"`
}

func handleGetLemmasMatchingText(body []byte) (interface{}, error) {
	var req getLemmasMatchingTextRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	if err := wordsmith.WithWordsmithTx(func(tx *sqlx.Tx) error {
		return nil
	}); err != nil {
		return nil, err
	}
	return nil, nil
}
