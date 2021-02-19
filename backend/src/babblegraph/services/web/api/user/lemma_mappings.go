package user

import (
	"babblegraph/model/routes"
	"babblegraph/model/userlemma"
	"babblegraph/services/web/util/routetoken"
	"babblegraph/util/database"
	"babblegraph/wordsmith"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type getUserLemmasForTokenRequest struct {
	Token string `json:"token"`
}

type getUserLemmasForTokenResponse struct {
	LemmaMappingsByLanguageCode []lemmaMappingsWithLanguageCode `json:"lemma_mappings_by_language_code"`
}

type lemmaMappingsWithLanguageCode struct {
	LanguageCode  wordsmith.LanguageCode `json:"language_code"`
	LemmaMappings []userlemma.Mapping    `json:"lemma_mappings"`
}

func handleGetUserLemmasForToken(body []byte) (interface{}, error) {
	var req getUserContentTopicsForTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.Token, routes.WordReinforcementKey)
	if err != nil {
		return nil, err
	}
	var lemmaMappings []userlemma.Mapping
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		lemmaMappings, err = userlemma.GetActiveMappingsForUser(tx, *userID)
		return err
	}); err != nil {
		return nil, err
	}
	lemmaMappingsByLanguageCode := make(map[wordsmith.LanguageCode][]userlemma.Mapping)
	for _, mapping := range lemmaMappings {
		mappings, _ := lemmaMappingsByLanguageCode[mapping.LanguageCode]
		lemmaMappingsByLanguageCode[mapping.LanguageCode] = append(mappings, mapping)
	}
	var out []lemmaMappingsWithLanguageCode
	for languageCode, mappings := range lemmaMappingsByLanguageCode {
		out = append(out, lemmaMappingsWithLanguageCode{
			LanguageCode:  languageCode,
			LemmaMappings: mappings,
		})
	}
	return getUserLemmasForTokenResponse{
		LemmaMappingsByLanguageCode: out,
	}, nil
}

type addUserLemmasForTokenRequest struct {
	Token   string            `json:"token"`
	LemmaID wordsmith.LemmaID `json:"lemma_id"`
}

type addUserLemmasForTokenResponse struct {
	DidUpdate bool `json:"did_update"`
}

func handleAddUserLemmasForToken(body []byte) (interface{}, error) {
	var req addUserLemmasForTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.Token, routes.WordReinforcementKey)
	if err != nil {
		return nil, err
	}
	var languageCode wordsmith.LanguageCode
	if err := wordsmith.WithWordsmithTx(func(tx *sqlx.Tx) error {
		lemma, err := wordsmith.GetLemmaByID(tx, req.LemmaID)
		if err != nil {
			return err
		}
		languageCode = lemma.Language
		return nil
	}); err != nil {
		return nil, err
	}
	var didUpdate bool
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		didUpdate, err = userlemma.AddMappingForUser(tx, *userID, req.LemmaID, languageCode)
		return err
	}); err != nil {
		return nil, err
	}
	return addUserLemmasForTokenResponse{
		DidUpdate: didUpdate,
	}, nil
}

type setUserLemmasInactiveForTokenRequest struct {
	Token     string              `json:"token"`
	MappingID userlemma.MappingID `json:"mapping_id"`
}

type setUserLemmasInactiveForTokenResponse struct {
	MappingID userlemma.MappingID `json:"mapping_id"`
	DidUpdate bool                `json:"did_update"`
}

func handleSetUserLemmasInactiveForToken(body []byte) (interface{}, error) {
	var req setUserLemmasInactiveForTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.Token, routes.WordReinforcementKey)
	if err != nil {
		return nil, err
	}
	var didUpdate bool
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		didUpdate, err = userlemma.SetMappingAsNotVisible(tx, *userID, req.MappingID)
		return err
	}); err != nil {
		return nil, err
	}
	return setUserLemmasInactiveForTokenResponse{
		MappingID: req.MappingID,
		DidUpdate: didUpdate,
	}, nil
}
