package contenttopics

import "fmt"

func GetContentTopicForString(topicString string) (*ContentTopic, error) {
	switch topicString {
	case "art":
		return ContentTopicArt.Ptr(), nil
	case "astronomy":
		return ContentTopicAstronomy.Ptr(), nil
	case "architecture":
		return ContentTopicArchitecture.Ptr(), nil
	case "automotive":
		return ContentTopicAutomotive.Ptr(), nil
	case "business":
		return ContentTopicBusiness.Ptr(), nil
	case "celebrity-news":
		return ContentTopicCelebrityNews.Ptr(), nil
	case "cooking":
		return ContentTopicCooking.Ptr(), nil
	case "culture":
		return ContentTopicCulture.Ptr(), nil
	case "current-events-argentina":
		return ContentTopicCurrentEventsArgentina.Ptr(), nil
	case "current-events-chile":
		return ContentTopicCurrentEventsChile.Ptr(), nil
	case "current-events-colombia":
		return ContentTopicCurrentEventsColombia.Ptr(), nil
	case "current-events-costa-rica":
		return ContentTopicCurrentEventsCostaRica.Ptr(), nil
	case "current-events-el-salvador":
		return ContentTopicCurrentEventsElSalvador.Ptr(), nil
	case "current-events-guatemala":
		return ContentTopicCurrentEventsGuatemala.Ptr(), nil
	case "current-events-honduras":
		return ContentTopicCurrentEventsHonduras.Ptr(), nil
	case "current-events-mexico":
		return ContentTopicCurrentEventsMexico.Ptr(), nil
	case "current-events-nicaragua":
		return ContentTopicCurrentEventsNicaragua.Ptr(), nil
	case "current-events-panama":
		return ContentTopicCurrentEventsPanama.Ptr(), nil
	case "current-events-paraguay":
		return ContentTopicCurrentEventsParaguay.Ptr(), nil
	case "current-events-peru":
		return ContentTopicCurrentEventsPeru.Ptr(), nil
	case "current-events-spain":
		return ContentTopicCurrentEventsSpain.Ptr(), nil
	case "current-events-united-states":
		return ContentTopicCurrentEventsUnitedStates.Ptr(), nil
	case "current-events-venezuela":
		return ContentTopicCurrentEventsVenezuela.Ptr(), nil
	case "current-events-uruguay":
		return ContentTopicCurrentEventsUruguay.Ptr(), nil
	case "current-events-puerto-rico":
		return ContentTopicCurrentEventsPuertoRico.Ptr(), nil
	case "economy":
		return ContentTopicEconomy.Ptr(), nil
	case "entertainment":
		return ContentTopicEntertainment.Ptr(), nil
	case "environment":
		return ContentTopicEnvironment.Ptr(), nil
	case "fashion":
		return ContentTopicFashion.Ptr(), nil
	case "film":
		return ContentTopicFilm.Ptr(), nil
	case "finance":
		return ContentTopicFinance.Ptr(), nil
	case "health":
		return ContentTopicHealth.Ptr(), nil
	case "home":
		return ContentTopicHome.Ptr(), nil
	case "lifestyle":
		return ContentTopicLifestyle.Ptr(), nil
	case "literature":
		return ContentTopicLiterature.Ptr(), nil
	case "music":
		return ContentTopicMusic.Ptr(), nil
	case "opinion":
		return ContentTopicOpinion.Ptr(), nil
	case "politics":
		return ContentTopicPolitics.Ptr(), nil
	case "science":
		return ContentTopicScience.Ptr(), nil
	case "sports":
		return ContentTopicSports.Ptr(), nil
	case "technology":
		return ContentTopicTechnology.Ptr(), nil
	case "theater":
		return ContentTopicTheater.Ptr(), nil
	case "travel":
		return ContentTopicTravel.Ptr(), nil
	case "video-games":
		return ContentTopicVideoGames.Ptr(), nil
	case "world-news":
		return ContentTopicWorldNews.Ptr(), nil
	}
	return nil, fmt.Errorf("unrecognized content topic: %s", topicString)
}
