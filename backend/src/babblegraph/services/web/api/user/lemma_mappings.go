package user

import (
	"babblegraph/model/routes"
	"babblegraph/model/userlemma"
	language_model "babblegraph/services/web/model/language"
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
	LemmaMappings []lemmaMapping `json:"lemma_mappings"`
}

type lemmaMapping struct {
	Lemma    language_model.Lemma `json:"lemma"`
	IsActive bool                 `json:"is_active"`
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
		lemmaMappings, err = userlemma.GetVisibleMappingsForUser(tx, *userID)
		return err
	}); err != nil {
		return nil, err
	}
	var lemmaIDs []wordsmith.LemmaID
	activeStatusByLemmaID := make(map[wordsmith.LemmaID]bool)
	for _, l := range lemmaMappings {
		lemmaIDs = append(lemmaIDs, l.LemmaID)
		activeStatusByLemmaID[l.LemmaID] = l.IsActive
	}
	lemmas, err := language_model.GetLemmasForLemmaIDs(lemmaIDs)
	if err != nil {
		return nil, err
	}
	var out []lemmaMapping
	for _, l := range lemmas {
		out = append(out, lemmaMapping{
			IsActive: activeStatusByLemmaID[l.ID],
			Lemma:    l,
		})
	}
	return getUserLemmasForTokenResponse{
		LemmaMappings: out,
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

type updateUserLemmaActiveStateForTokenRequest struct {
	Token        string            `json:"token"`
	LemmaID      wordsmith.LemmaID `json:"lemma_id"`
	CurrentState bool              `json:"current_state"`
}

type updateUserLemmaActiveStateForTokenResponse struct {
	LemmaID   wordsmith.LemmaID `json:"lemma_id"`
	DidUpdate bool              `json:"did_update"`
}

func handleUpdateUserLemmaActiveStateForToken(body []byte) (interface{}, error) {
	var req updateUserLemmaActiveStateForTokenRequest
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
		didUpdate, err = userlemma.ToggleMappingActiveState(tx, *userID, req.LemmaID, req.CurrentState)
		return err
	}); err != nil {
		return nil, err
	}
	return updateUserLemmaActiveStateForTokenResponse{
		LemmaID:   req.LemmaID,
		DidUpdate: didUpdate,
	}, nil
}

type removeUserLemmaForTokenRequest struct {
	Token   string            `json:"token"`
	LemmaID wordsmith.LemmaID `json:"lemma_id"`
}

type removeUserLemmaForTokenResponse struct {
	LemmaID   wordsmith.LemmaID `json:"lemma_id"`
	DidDelete bool              `json:"did_delete"`
}

func removeUserLemmaForToken(body []byte) (interface{}, error) {
	var req removeUserLemmaForTokenRequest
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
	var didDelete bool
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		didDelete, err = userlemma.RemoveMappingForUser(tx, *userID, req.LemmaID, languageCode)
		return err
	}); err != nil {
		return nil, err
	}
	return removeUserLemmaForTokenResponse{
		LemmaID:   req.LemmaID,
		DidDelete: didDelete,
	}, nil
}