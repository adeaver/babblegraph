package content

import (
	"babblegraph/model/content"
	"babblegraph/services/web/clientrouter/clienterror"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "content",
	Routes: []router.Route{
		{
			Path: "get_active_topics_for_language_code_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				getActiveTopicsForLanguageCode,
			),
		},
	},
}

type getActiveTopicsForLanguageCodeRequest struct {
	LanguageCode string `json:"language_code"`
}

type topicWithDisplay struct {
	Topic       content.Topic            `json:"topic"`
	DisplayName content.TopicDisplayName `json:"topic_display_name"`
}

type getActiveTopicsForLanguageCodeResponse struct {
	Topics []topicWithDisplay `json:"topics,omitempty"`
	Error  *clienterror.Error `json:"error,omitempty"`
}

func getActiveTopicsForLanguageCode(r *router.Request) (interface{}, error) {
	var req getActiveTopicsForLanguageCodeRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return getActiveTopicsForLanguageCodeResponse{
			Error: clienterror.ErrorInvalidLanguageCode.Ptr(),
		}, nil
	}
	var topics []content.Topic
	var displayNames []content.TopicDisplayName
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		topics, err = content.GetAllTopics(tx)
		if err != nil {
			return err
		}
		displayNames, err = content.GetTopicDisplayNamesForLanguage(tx, *languageCode)
		return err
	}); err != nil {
		return nil, err
	}
	topicDisplayNamesByTopicID := make(map[content.TopicID]content.TopicDisplayName)
	for _, d := range displayNames {
		topicDisplayNamesByTopicID[d.TopicID] = d
	}
	var out []topicWithDisplay
	for _, t := range topics {
		if !t.IsActive {
			continue
		}
		displayName, ok := topicDisplayNamesByTopicID[t.ID]
		if !ok {
			continue
		}
		out = append(out, topicWithDisplay{
			Topic:       t,
			DisplayName: displayName,
		})
	}
	return getActiveTopicsForLanguageCodeResponse{
		Topics: out,
	}, nil
}
