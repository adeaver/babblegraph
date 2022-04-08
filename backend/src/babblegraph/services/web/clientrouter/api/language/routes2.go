package language

import (
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/services/web/clientrouter/clienterror"
	language_model "babblegraph/services/web/clientrouter/model/language"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/clientrouter/util/routetoken"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/recaptcha"
	"babblegraph/wordsmith"
	"fmt"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "language",
	Routes: []router.Route{
		{
			Path: "search_text_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				routermiddleware.MaybeWithAuthentication(searchText),
			),
		},
	},
}

type searchTextRequest struct {
	CaptchaToken           string   `json:"captcha_token"`
	WordReinforcementToken string   `json:"word_reinforcement_token"`
	Text                   []string `json:"text"`
	LanguageCode           string   `json:"language_code"`
}

type searchTextResponse struct {
	Error  *clienterror.Error `json:"error,omitempty"`
	Result *searchResult      `json:"result,omitempty"`
}

type searchResult struct {
	Results      []textSearchResult     `json:"results,omitempty"`
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
}

// TODO: move this
type textSearchResult struct {
	DisplayText  string                       `json:"display_text"`
	Definitions  []string                     `json:"definitions,omitempty"`
	PartOfSpeech *language_model.PartOfSpeech `json:"part_of_speech,omitempty"`
	LookupID     textSearchLookupID           `json:"lookup_id"`
}

type textSearchDefinition struct {
	EnglishDefinitionDisplay []string     `json:"english_definition_display"`
	DefinitionID             definitionID `json:"definition_id"`
}

type definitionID struct {
	IDType textSearchLookupIDType `json:"id_type"`
	ID     string                 `json:"id"`
}

type textSearchLookupID struct {
	IDType textSearchLookupIDType `json:"id_type"`
	ID     []string               `json:"id"`
}

type textSearchLookupIDType string

const (
	textSearchLookupIDTypeLemma  textSearchLookupIDType = "lemma"
	textSearchLookupIDTypePhrase textSearchLookupIDType = "phrase"
)

const (
	errorInvalidSearchLength clienterror.Error = "invalid-search-length"
)

func searchText(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req *searchTextRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.WordReinforcementToken, routes.WordReinforcementKey)
	if err != nil {
		return searchTextResponse{
			Error: clienterror.ErrorInvalidToken.Ptr(),
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
		return searchTextResponse{
			Error: clienterror.ErrorIncorrectKey.Ptr(),
		}, nil
	case userAuth == nil && doesUserHaveAccount:
		r.RespondWithStatus(http.StatusForbidden)
		return searchTextResponse{
			Error: clienterror.ErrorNoAuth.Ptr(),
		}, nil
	case userAuth.SubscriptionLevel == nil && len(req.Text) > 1:
		return searchTextResponse{
			Error: clienterror.ErrorRequiresUpgrade.Ptr(),
		}, nil
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return searchTextResponse{
			Error: clienterror.ErrorInvalidLanguageCode.Ptr(),
		}, nil
	}
	isValid, err := recaptcha.VerifyRecaptchaToken("searchtext", req.CaptchaToken)
	switch {
	case err != nil:
		return nil, err
	case !isValid:
		return searchTextResponse{
			Error: clienterror.ErrorLowCaptchaScore.Ptr(),
		}, nil
	default:
		r.Infof("Successfully cleared captcha")
	}
	switch {
	case len(req.Text) == 1:
		lemmas, err := language_model.GetLemmasForWordText(req.Text[0])
		if err != nil {
			return nil, err
		}
		var result []textSearchResult
		for _, l := range lemmas {
			l := l
			var definitions []string
			for _, d := range l.Definitions {
				definitions = append(definitions, d.Text)
			}
			result = append(result, textSearchResult{
				DisplayText:  l.Text,
				Definitions:  definitions,
				PartOfSpeech: &l.PartOfSpeech,
				LookupID: textSearchLookupID{
					IDType: textSearchLookupIDTypeLemma,
					ID:     []string{l.ID.Str()},
				},
			})
		}
		return searchTextResponse{
			Result: &searchResult{
				Results:      result,
				LanguageCode: *languageCode,
			},
		}, nil
	case len(req.Text) > 1:
		var phraseDefinitions []wordsmith.PhraseDefinition
		foundAll := true
		if err := wordsmith.WithWordsmithTx(func(tx *sqlx.Tx) error {
			words, err := wordsmith.GetWordsByText(tx, wordsmith.SpanishUPCWikiCorpus, req.Text)
			if err != nil {
				return err
			}
			lemmaIDsByText := make(map[string][]wordsmith.LemmaID)
			var lemmaIDs []wordsmith.LemmaID
			for _, w := range words {
				lemmaIDs = append(lemmaIDs, w.LemmaID)
				lemmaIDsByText[strings.ToLower(w.WordText)] = append(lemmaIDsByText[strings.ToLower(w.WordText)], w.LemmaID)
			}
			lemmas, err := wordsmith.GetLemmasByIDs(tx, lemmaIDs)
			if err != nil {
				return err
			}
			lemmaTextByID := make(map[wordsmith.LemmaID]string)
			for _, lemma := range lemmas {
				lemmaTextByID[lemma.ID] = lemma.LemmaText
			}
			var lemmaTexts [][]string
			for idx, w := range req.Text {
				lemmaTexts = append(lemmaTexts, make([]string, 0))
				lemmaIDs, _ := lemmaIDsByText[strings.ToLower(w)]
				for _, lemmaID := range lemmaIDs {
					lemmaText, ok := lemmaTextByID[lemmaID]
					if !ok {
						continue
					}
					lemmaTexts[idx] = append(lemmaTexts[idx], lemmaText)
				}
				if len(lemmaTexts[idx]) == 0 {
					foundAll = false
					return nil
				}
			}
			r.Debugf("Got lemmas %+v", lemmaTexts)
			lemmaPhrases := makeLemmaPhrases(lemmaTexts, "")
			r.Debugf("Got phrases %+v (length %d)", lemmaPhrases, len(lemmaPhrases))
			phraseDefinitions, err = wordsmith.GetPhraseDefinitionsForLemmaPhrases(tx, wordsmith.SpanishUPCWikiCorpus, lemmaPhrases)
			return err
		}); err != nil {
			return nil, err
		}
		if !foundAll {
			return searchTextResponse{
				Result: &searchResult{
					Results:      []textSearchResult{},
					LanguageCode: *languageCode,
				},
			}, nil
		}
		var result []textSearchResult
		for _, p := range phraseDefinitions {
			result = append(result, textSearchResult{
				DisplayText: p.Phrase,
				Definitions: []string{p.Definition},
				LookupID: textSearchLookupID{
					IDType: textSearchLookupIDTypePhrase,
					ID:     []string{string(p.ID)},
				},
			})
		}
		result = append(result, textSearchResult{
			DisplayText: strings.Join(req.Text, " "),
			LookupID: textSearchLookupID{
				IDType: textSearchLookupIDTypePhrase,
				ID:     req.Text,
			},
		})
		return searchTextResponse{
			Result: &searchResult{
				Results:      result,
				LanguageCode: *languageCode,
			},
		}, nil
	default:
		return searchTextResponse{
			Error: errorInvalidSearchLength.Ptr(),
		}, nil
	}
}

func makeLemmaPhrases(lemmaTexts [][]string, currentPhrase string) []string {
	if len(lemmaTexts) == 0 {
		return []string{strings.Trim(currentPhrase, " ")}
	}
	var out []string
	for _, lemma := range lemmaTexts[0] {
		out = append(out, makeLemmaPhrases(lemmaTexts[1:], fmt.Sprintf("%s %s", currentPhrase, lemma))...)
	}
	return out
}
