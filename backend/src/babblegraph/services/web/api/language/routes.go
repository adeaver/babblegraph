package language

import (
	"babblegraph/model/routes"
	language_model "babblegraph/services/web/model/language"
	"babblegraph/services/web/router"
	"babblegraph/util/encrypt"
	"babblegraph/wordsmith"
	"encoding/json"
	"fmt"
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
	Lemmas       []language_model.Lemma `json:"lemmas"`
}

func handleGetLemmasMatchingText(body []byte) (interface{}, error) {
	var req getLemmasMatchingTextRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	if err := encrypt.WithDecodedToken(req.Token, func(t encrypt.TokenPair) error {
		if t.Key != routes.WordReinforcementKey.Str() {
			return fmt.Errorf("incorrect key")
		}
		return nil
	}); err != nil {
		return nil, err
	}
	lemmas, err := language_model.GetLemmasForWordText(req.Text)
	if err != nil {
		return nil, err
	}
	return getLemmasMatchingTextResponse{
		LanguageCode: req.LanguageCode,
		Text:         req.Text,
		Lemmas:       lemmas,
	}, nil
}
