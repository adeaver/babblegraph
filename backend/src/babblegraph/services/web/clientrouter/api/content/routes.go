package content

import (
	"babblegraph/model/content"
	"babblegraph/services/web/clientrouter/clienterror"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/wordsmith"
	"strings"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "content",
	Routes: []router.Route{
		{
			Path: "get_topics_for_language_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				routermiddleware.MaybeWithAuthentication(getTopicsForLanguage),
			),
		},
	},
}

type getTopicsForLanguageRequest struct {
	LanguageCode string `json:"language_code"`
}

type getTopicsForLanguageResponse struct {
	Error   *clienterror.Error `json:"error,omitempty"`
	Results []topic            `json:"results,omitempty"`
}

type topic struct {
	TopicID      content.TopicID `json:"topic_id"`
	DisplayName  string          `json:"display_name"`
	EnglishLabel string          `json:"english_label"`
}

func getTopicsForLanguage(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req *getTopicsForLanguageRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return getTopicsForLanguageResponse{
			Error: clienterror.ErrorInvalidLanguageCode.Ptr(),
		}, nil
	}
	var topics []content.Topic
	var topicDisplayNames []content.TopicDisplayName
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		topicDisplayNames, err = content.GetTopicDisplayNamesForLanguage(tx, *languageCode)
		if err != nil {
			return err
		}
		topics, err = content.GetAllTopics(tx)
		return err
	}); err != nil {
		return nil, err
	}
	topicsByID := make(map[content.TopicID]content.Topic)
	for _, t := range topics {
		topicsByID[t.ID] = t
	}
	var out []topic
	for _, d := range topicDisplayNames {
		t, ok := topicsByID[d.TopicID]
		if !d.IsActive || !ok || !t.IsActive {
			continue
		}
		readableLabel := strings.ReplaceAll(strings.TrimPrefix(t.Label, "current-events-"), "-", " ")
		out = append(out, topic{
			TopicID:      d.TopicID,
			DisplayName:  d.Label,
			EnglishLabel: readableLabel,
		})
	}
	return getTopicsForLanguageResponse{
		Results: out,
	}, nil
}
