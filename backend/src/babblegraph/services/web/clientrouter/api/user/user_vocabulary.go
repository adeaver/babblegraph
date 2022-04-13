package user

import (
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/uservocabulary"
	"babblegraph/services/web/clientrouter/clienterror"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/clientrouter/util/routetoken"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

type upsertUserVocabularyRequest struct {
	SubscriptionManagementToken string  `json:"subscription_management_token"`
	LanguageCode                string  `json:"language_code"`
	DisplayText                 string  `json:"display_text"`
	DefinitionID                *string `json:"definition_id,omitempty"`
	EntryType                   string  `json:"entry_type"`
	StudyNote                   *string `json:"study_note,omitempty"`
	IsVisible                   bool    `json:"is_visible"`
	IsActive                    bool    `json:"is_active"`
}

const (
	errorInvalidEntryType clienterror.Error = "invalid-entry-type"
	errorInvalidInput     clienterror.Error = "invalid-input"
)

type upsertUserVocabularyResponse struct {
	ID    *uservocabulary.UserVocabularyEntryID `json:"id,omitempty"`
	Error *clienterror.Error                    `json:"error"`
}

func upsertUserVocabulary(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req upsertUserVocabularyRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.SubscriptionManagementToken, routes.SubscriptionManagementRouteEncryptionKey)
	if err != nil {
		return upsertUserVocabularyResponse{
			Error: clienterror.ErrorInvalidToken.Ptr(),
		}, nil
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return upsertUserVocabularyResponse{
			Error: clienterror.ErrorInvalidLanguageCode.Ptr(),
		}, nil
	}
	entryType, err := uservocabulary.GetVocabularyTypeFromString(req.EntryType)
	if err != nil {
		return upsertUserVocabularyResponse{
			Error: errorInvalidEntryType.Ptr(),
		}, nil
	}
	var doesUserHaveAccount bool
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		doesUserHaveAccount, err = useraccounts.DoesUserAlreadyHaveAccount(tx, *userID)
		return err
	}); err != nil {
		return nil, err
	}
	switch {
	case userAuth != nil && userAuth.UserID != *userID:
		return upsertUserVocabularyResponse{
			Error: clienterror.ErrorIncorrectKey.Ptr(),
		}, nil
	case userAuth == nil && doesUserHaveAccount:
		return upsertUserVocabularyResponse{
			Error: clienterror.ErrorNoAuth.Ptr(),
		}, nil
	case userAuth == nil && *entryType == uservocabulary.VocabularyTypePhrase:
		return upsertUserVocabularyResponse{
			Error: clienterror.ErrorRequiresUpgrade.Ptr(),
		}, nil
	}
	input := uservocabulary.UpsertVocabularyEntryInput{
		UserID:       *userID,
		LanguageCode: *languageCode,
		IsActive:     req.IsActive,
		IsVisible:    req.IsVisible,
		StudyNote:    req.StudyNote,
	}
	switch *entryType {
	case uservocabulary.VocabularyTypePhrase:
		hashablePhrase := &uservocabulary.HashablePhrase{
			DisplayText: req.DisplayText,
		}
		if req.DefinitionID != nil {
			phraseDefinitionID := wordsmith.PhraseDefinitionID(*req.DefinitionID)
			hashablePhrase.DefinitionID = &phraseDefinitionID
		}
		input.Hashable = hashablePhrase
	case uservocabulary.VocabularyTypeLemma:
		if req.DefinitionID == nil {
			return upsertUserVocabularyResponse{
				Error: errorInvalidInput.Ptr(),
			}, nil
		}
		input.Hashable = &uservocabulary.HashableLemma{
			LemmaID:   wordsmith.LemmaID(*req.DefinitionID),
			LemmaText: req.DisplayText,
		}
	default:
		return upsertUserVocabularyResponse{
			Error: errorInvalidInput.Ptr(),
		}, nil
	}
	var userVocabularyID *uservocabulary.UserVocabularyEntryID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		userVocabularyID, err = uservocabulary.UpsertVocabularyEntry(tx, input)
		return err
	}); err != nil {
		return nil, err
	}
	return upsertUserVocabularyResponse{
		ID: userVocabularyID,
	}, nil
}

type getUserVocabularyRequest struct {
	SubscriptionManagementToken string `json:"subscription_management_token"`
	LanguageCode                string `json:"language_code"`
}

type getUserVocabularyResponse struct {
	Entries []uservocabulary.UserVocabularyEntry `json:"entries,omitempty"`
	Error   *clienterror.Error                   `json:"error,omitempty"`
}

func getUserVocabulary(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req getUserVocabularyRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.SubscriptionManagementToken, routes.SubscriptionManagementRouteEncryptionKey)
	if err != nil {
		return getUserVocabularyResponse{
			Error: clienterror.ErrorInvalidToken.Ptr(),
		}, nil
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return getUserVocabularyResponse{
			Error: clienterror.ErrorInvalidLanguageCode.Ptr(),
		}, nil
	}
	var doesUserHaveAccount bool
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		doesUserHaveAccount, err = useraccounts.DoesUserAlreadyHaveAccount(tx, *userID)
		return err
	}); err != nil {
		return nil, err
	}
	switch {
	case userAuth != nil && userAuth.UserID != *userID:
		return getUserVocabularyResponse{
			Error: clienterror.ErrorIncorrectKey.Ptr(),
		}, nil
	case userAuth == nil && doesUserHaveAccount:
		return getUserVocabularyResponse{
			Error: clienterror.ErrorNoAuth.Ptr(),
		}, nil
	}
	var entries []uservocabulary.UserVocabularyEntry
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		entries, err = uservocabulary.GetUserVocabularyEntries(tx, *userID, *languageCode, true)
		return err
	}); err != nil {
		return nil, err
	}
	return getUserVocabularyResponse{
		Entries: entries,
	}, nil
}
