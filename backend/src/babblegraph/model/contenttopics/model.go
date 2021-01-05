package contenttopics

type ContentTopic string

const (
	ContentTopicEconomy                ContentTopic = "economy"
	ContentTopicFinance                ContentTopic = "finance"
	ContentTopicPolitics               ContentTopic = "politics"
	ContentTopicBusiness               ContentTopic = "business"
	ContentTopicLifestyle              ContentTopic = "lifestyle"
	ContentTopicOpinion                ContentTopic = "opinion"
	ContentTopicCurrentEventsArgentina ContentTopic = "current-events-argentina"
	ContentTopicWorldNews              ContentTopic = "world-news"
	ContentTopicEntertainment          ContentTopic = "entertainment"
	ContentTopicAutomotive             ContentTopic = "automotive"
	ContentTopicSports                 ContentTopic = "sports"
)

type contentTopicMappingID string

type dbContentTopicMapping struct {
	ID            contentTopicMappingID `db:"_id"`
	URLIdentifier string                `db:"url_identifier"`
	ContentTopic  ContentTopic          `db:"content_topic"`
}
