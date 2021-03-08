package contenttopics

import "fmt"

type DisplayName string

const (
	DisplayNameArt                       DisplayName = "arte"
	DisplayNameArchitecture              DisplayName = "arquitectura"
	DisplayNameAutomotive                DisplayName = "autos"
	DisplayNameBusiness                  DisplayName = "negocios"
	DisplayNameCelebrityNews             DisplayName = "noticias de celebridades"
	DisplayNameCooking                   DisplayName = "cocina"
	DisplayNameCulture                   DisplayName = "cultura"
	DisplayNameCurrentEventsArgentina    DisplayName = "noticias de argentina"
	DisplayNameCurrentEventsChile        DisplayName = "noticias de chile"
	DisplayNameCurrentEventsColombia     DisplayName = "noticias de colombia"
	DisplayNameCurrentEventsCostaRica    DisplayName = "noticias de costa rica"
	DisplayNameCurrentEventsElSalvador   DisplayName = "noticias de el salvador"
	DisplayNameCurrentEventsGuatemala    DisplayName = "noticias de guatemala"
	DisplayNameCurrentEventsHonduras     DisplayName = "noticias de honduras"
	DisplayNameCurrentEventsMexico       DisplayName = "noticias de méxico"
	DisplayNameCurrentEventsNicaragua    DisplayName = "noticias de nicaragua"
	DisplayNameCurrentEventsPanama       DisplayName = "noticias de panamá"
	DisplayNameCurrentEventsParaguay     DisplayName = "noticias de paraguay"
	DisplayNameCurrentEventsPeru         DisplayName = "noticias de perú"
	DisplayNameCurrentEventsSpain        DisplayName = "noticias de españa"
	DisplayNameCurrentEventsUnitedStates DisplayName = "noticias de estados unidos"
	DisplayNameCurrentEventsVenezuela    DisplayName = "noticias de venezuela"
	DisplayNameCurrentEventsUruguay      DisplayName = "noticias de uruguay"
	DisplayNameEconomy                   DisplayName = "economía"
	DisplayNameEntertainment             DisplayName = "entretenimiento"
	DisplayNameFashion                   DisplayName = "moda"
	DisplayNameFilm                      DisplayName = "cine"
	DisplayNameFinance                   DisplayName = "finanzas"
	DisplayNameHealth                    DisplayName = "salud"
	DisplayNameHome                      DisplayName = "casa"
	DisplayNameLifestyle                 DisplayName = "estilo de vida"
	DisplayNameLiterature                DisplayName = "literatura"
	DisplayNameMusic                     DisplayName = "música"
	DisplayNameOpinion                   DisplayName = "opinión"
	DisplayNamePolitics                  DisplayName = "política"
	DisplayNameScience                   DisplayName = "ciencia"
	DisplayNameSports                    DisplayName = "deportes"
	DisplayNameTechnology                DisplayName = "tecnología"
	DisplayNameTheater                   DisplayName = "teatro"
	DisplayNameTravel                    DisplayName = "viajes"
	DisplayNameVideoGames                DisplayName = "videojuegos"
	DisplayNameWorldNews                 DisplayName = "noticias del mundo"
)

func (d DisplayName) Str() string {
	return string(d)
}

func (d DisplayName) Ptr() *DisplayName {
	return &d
}

func ContentTopicNameToDisplayName(topic ContentTopic) (*DisplayName, error) {
	switch topic {
	case ContentTopicArt:
		return DisplayNameArt.Ptr(), nil
	case ContentTopicArchitecture:
		return DisplayNameArchitecture.Ptr(), nil
	case ContentTopicAutomotive:
		return DisplayNameAutomotive.Ptr(), nil
	case ContentTopicBusiness:
		return DisplayNameBusiness.Ptr(), nil
	case ContentTopicCelebrityNews:
		return DisplayNameCelebrityNews.Ptr(), nil
	case ContentTopicCooking:
		return DisplayNameCooking.Ptr(), nil
	case ContentTopicCulture:
		return DisplayNameCulture.Ptr(), nil
	case ContentTopicCurrentEventsArgentina:
		return DisplayNameCurrentEventsArgentina.Ptr(), nil
	case ContentTopicCurrentEventsChile:
		return DisplayNameCurrentEventsChile.Ptr(), nil
	case ContentTopicCurrentEventsColombia:
		return DisplayNameCurrentEventsColombia.Ptr(), nil
	case ContentTopicCurrentEventsCostaRica:
		return DisplayNameCurrentEventsCostaRica.Ptr(), nil
	case ContentTopicCurrentEventsElSalvador:
		return DisplayNameCurrentEventsElSalvador.Ptr(), nil
	case ContentTopicCurrentEventsGuatemala:
		return DisplayNameCurrentEventsGuatemala.Ptr(), nil
	case ContentTopicCurrentEventsHonduras:
		return DisplayNameCurrentEventsHonduras.Ptr(), nil
	case ContentTopicCurrentEventsMexico:
		return DisplayNameCurrentEventsMexico.Ptr(), nil
	case ContentTopicCurrentEventsNicaragua:
		return DisplayNameCurrentEventsNicaragua.Ptr(), nil
	case ContentTopicCurrentEventsPanama:
		return DisplayNameCurrentEventsPanama.Ptr(), nil
	case ContentTopicCurrentEventsParaguay:
		return DisplayNameCurrentEventsParaguay.Ptr(), nil
	case ContentTopicCurrentEventsPeru:
		return DisplayNameCurrentEventsPeru.Ptr(), nil
	case ContentTopicCurrentEventsSpain:
		return DisplayNameCurrentEventsSpain.Ptr(), nil
	case ContentTopicCurrentEventsUnitedStates:
		return DisplayNameCurrentEventsUnitedStates.Ptr(), nil
	case ContentTopicCurrentEventsVenezuela:
		return DisplayNameCurrentEventsVenezuela.Ptr(), nil
	case ContentTopicCurrentEventsUruguay:
		return DisplayNameCurrentEventsUruguay.Ptr(), nil
	case ContentTopicEconomy:
		return DisplayNameEconomy.Ptr(), nil
	case ContentTopicEntertainment:
		return DisplayNameEntertainment.Ptr(), nil
	case ContentTopicFashion:
		return DisplayNameFashion.Ptr(), nil
	case ContentTopicFilm:
		return DisplayNameFilm.Ptr(), nil
	case ContentTopicFinance:
		return DisplayNameFinance.Ptr(), nil
	case ContentTopicHealth:
		return DisplayNameHealth.Ptr(), nil
	case ContentTopicHome:
		return DisplayNameHome.Ptr(), nil
	case ContentTopicLifestyle:
		return DisplayNameLifestyle.Ptr(), nil
	case ContentTopicLiterature:
		return DisplayNameLiterature.Ptr(), nil
	case ContentTopicMusic:
		return DisplayNameMusic.Ptr(), nil
	case ContentTopicOpinion:
		return DisplayNameOpinion.Ptr(), nil
	case ContentTopicPolitics:
		return DisplayNamePolitics.Ptr(), nil
	case ContentTopicScience:
		return DisplayNameScience.Ptr(), nil
	case ContentTopicSports:
		return DisplayNameSports.Ptr(), nil
	case ContentTopicTechnology:
		return DisplayNameTechnology.Ptr(), nil
	case ContentTopicTheater:
		return DisplayNameTheater.Ptr(), nil
	case ContentTopicTravel:
		return DisplayNameTravel.Ptr(), nil
	case ContentTopicVideoGames:
		return DisplayNameVideoGames.Ptr(), nil
	case ContentTopicWorldNews:
		return DisplayNameWorldNews.Ptr(), nil
	default:
		return nil, fmt.Errorf("unsupported topic: %s", topic.Str())
	}
}
