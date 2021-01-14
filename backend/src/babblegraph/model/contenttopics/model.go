package contenttopics

type ContentTopic string

const (
	ContentTopicArt                       ContentTopic = "art"
	ContentTopicArchitecture              ContentTopic = "architecture"
	ContentTopicAutomotive                ContentTopic = "automotive"
	ContentTopicBusiness                  ContentTopic = "business"
	ContentTopicCelebrityNews             ContentTopic = "celebrity-news"
	ContentTopicCooking                   ContentTopic = "cooking"
	ContentTopicCulture                   ContentTopic = "culture"
	ContentTopicCurrentEventsArgentina    ContentTopic = "current-events-argentina"
	ContentTopicCurrentEventsChile        ContentTopic = "current-events-chile"
	ContentTopicCurrentEventsColombia     ContentTopic = "current-events-colombia"
	ContentTopicCurrentEventsCostaRica    ContentTopic = "current-events-costa-rica"
	ContentTopicCurrentEventsElSalvador   ContentTopic = "current-events-el-salvador"
	ContentTopicCurrentEventsGuatemala    ContentTopic = "current-events-guatemala"
	ContentTopicCurrentEventsHonduras     ContentTopic = "current-events-honduras"
	ContentTopicCurrentEventsMexico       ContentTopic = "current-events-mexico"
	ContentTopicCurrentEventsNicaragua    ContentTopic = "current-events-nicaragua"
	ContentTopicCurrentEventsPanama       ContentTopic = "current-events-panama"
	ContentTopicCurrentEventsPeru         ContentTopic = "current-events-peru"
	ContentTopicCurrentEventsParaguay     ContentTopic = "current-events-paraguay"
	ContentTopicCurrentEventsSpain        ContentTopic = "current-events-spain"
	ContentTopicCurrentEventsUnitedStates ContentTopic = "current-events-united-states"
	ContentTopicCurrentEventsVenezuela    ContentTopic = "current-events-venezuela"
	ContentTopicCurrentEventsUruguay      ContentTopic = "current-events-uruguay"
	ContentTopicEconomy                   ContentTopic = "economy"
	ContentTopicEntertainment             ContentTopic = "entertainment"
	ContentTopicFashion                   ContentTopic = "fashion"
	ContentTopicFilm                      ContentTopic = "film"
	ContentTopicFinance                   ContentTopic = "finance"
	ContentTopicHealth                    ContentTopic = "health"
	ContentTopicHome                      ContentTopic = "home"
	ContentTopicLifestyle                 ContentTopic = "lifestyle"
	ContentTopicLiterature                ContentTopic = "literature"
	ContentTopicMusic                     ContentTopic = "music"
	ContentTopicOpinion                   ContentTopic = "opinion"
	ContentTopicPolitics                  ContentTopic = "politics"
	ContentTopicScience                   ContentTopic = "science"
	ContentTopicSports                    ContentTopic = "sports"
	ContentTopicTechnology                ContentTopic = "technology"
	ContentTopicTheater                   ContentTopic = "theater"
	ContentTopicTravel                    ContentTopic = "travel"
	ContentTopicVideoGames                ContentTopic = "video-games"
	ContentTopicWorldNews                 ContentTopic = "world-news"
)

func (t ContentTopic) Str() string {
	return string(t)
}

func (t ContentTopic) Ptr() *ContentTopic {
	return &t
}

type contentTopicMappingID string

type dbContentTopicMapping struct {
	ID            contentTopicMappingID `db:"_id"`
	URLIdentifier string                `db:"url_identifier"`
	ContentTopic  ContentTopic          `db:"content_topic"`
}
