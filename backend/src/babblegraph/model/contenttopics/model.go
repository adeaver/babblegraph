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
	ContentTopicCurrentEventsChile     ContentTopic = "current-events-chile"
	ContentTopicWorldNews              ContentTopic = "world-news"
	ContentTopicEntertainment          ContentTopic = "entertainment"
	ContentTopicAutomotive             ContentTopic = "automotive"
	ContentTopicSports                 ContentTopic = "sports"
	ContentTopicScience                ContentTopic = "science"
	ContentTopicHealth                 ContentTopic = "health"
	ContentTopicLiterature             ContentTopic = "literature"
	ContentTopicTechnology             ContentTopic = "technology"
	ContentTopicFashion                ContentTopic = "fashion"
	ContentTopicHome                   ContentTopic = "home"
	ContentTopicCooking                ContentTopic = "cooking"
	ContentTopicCelebrityNews          ContentTopic = "celebrity-news"
	ContentTopicMusic                  ContentTopic = "music"
	ContentTopicFilm                   ContentTopic = "film"
	ContentTopicVideoGames             ContentTopic = "video-games"
	ContentTopicArt                    ContentTopic = "art"
)

type contentTopicMappingID string

type dbContentTopicMapping struct {
	ID            contentTopicMappingID `db:"_id"`
	URLIdentifier string                `db:"url_identifier"`
	ContentTopic  ContentTopic          `db:"content_topic"`
}
