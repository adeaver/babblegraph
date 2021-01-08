package contenttopics

type ContentTopic string

const (
	ContentTopicArt                     ContentTopic = "art"
	ContentTopicAutomotive              ContentTopic = "automotive"
	ContentTopicBusiness                ContentTopic = "business"
	ContentTopicCelebrityNews           ContentTopic = "celebrity-news"
	ContentTopicCooking                 ContentTopic = "cooking"
	ContentTopicCulture                 ContentTopic = "culture"
	ContentTopicCurrentEventsArgentina  ContentTopic = "current-events-argentina"
	ContentTopicCurrentEventsChile      ContentTopic = "current-events-chile"
	ContentTopicCurrentEventsColombia   ContentTopic = "current-events-colombia"
	ContentTopicCurrentEventsCostaRica  ContentTopic = "current-events-costa-rica"
	ContentTopicCurrentEventsElSalvador ContentTopic = "current-events-el-salvador"
	ContentTopicEconomy                 ContentTopic = "economy"
	ContentTopicEntertainment           ContentTopic = "entertainment"
	ContentTopicFashion                 ContentTopic = "fashion"
	ContentTopicFilm                    ContentTopic = "film"
	ContentTopicFinance                 ContentTopic = "finance"
	ContentTopicHealth                  ContentTopic = "health"
	ContentTopicHome                    ContentTopic = "home"
	ContentTopicLifestyle               ContentTopic = "lifestyle"
	ContentTopicLiterature              ContentTopic = "literature"
	ContentTopicMusic                   ContentTopic = "music"
	ContentTopicOpinion                 ContentTopic = "opinion"
	ContentTopicPolitics                ContentTopic = "politics"
	ContentTopicScience                 ContentTopic = "science"
	ContentTopicSports                  ContentTopic = "sports"
	ContentTopicTechnology              ContentTopic = "technology"
	ContentTopicTravel                  ContentTopic = "travel"
	ContentTopicVideoGames              ContentTopic = "video-games"
	ContentTopicWorldNews               ContentTopic = "world-news"
)

type contentTopicMappingID string

type dbContentTopicMapping struct {
	ID            contentTopicMappingID `db:"_id"`
	URLIdentifier string                `db:"url_identifier"`
	ContentTopic  ContentTopic          `db:"content_topic"`
}
