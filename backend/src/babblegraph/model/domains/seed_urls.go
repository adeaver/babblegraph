package domains

import "babblegraph/model/contenttopics"

var seedURLs = []SeedURL{
	{
		URL: "https://www.ambito.com/contenidos/economia.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://www.ambito.com/contenidos/finanzas.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFinance,
		},
	}, {
		URL: "https://www.ambito.com/contenidos/politica.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicPolitics,
		},
	}, {
		URL: "https://www.ambito.com/contenidos/negocios.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicBusiness,
		},
	}, {
		URL: "https://www.ambito.com/contenidos/lifestyle.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicLifestyle,
		},
	}, {
		URL: "https://www.ambito.com/contenidos/opinion.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://www.ambito.com/contenidos/nacional.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsArgentina,
		},
	}, {
		URL: "https://www.ambito.com/contenidos/mundo.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.ambito.com/contenidos/espectaculos.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://www.ambito.com/contenidos/autos.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicAutomotive,
		},
	}, {
		URL: "https://www.ambito.com/contenidos/deportes.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://www.cronista.com/apertura-negocio",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicBusiness,
		},
	}, {
		URL: "https://www.cronista.com/seccion/economia_politica",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
			contenttopics.ContentTopicPolitics,
		},
	}, {
		URL: "https://www.cronista.com/seccion/finanzas_mercados",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFinance,
		},
	}, {
		URL: "https://www.cronista.com/seccion/internacionales",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.cronista.com/seccion/columnistas",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "http://www.laprensa.com.ar/category.aspx?category=114",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicPolitics,
		},
	}, {
		URL: "http://www.laprensa.com.ar/category.aspx?category=115",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "http://www.laprensa.com.ar/category.aspx?category=116",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "http://www.laprensa.com.ar/category.aspx?category=117",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "http://www.laprensa.com.ar/category.aspx?category=119",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "http://www.laprensa.com.ar/category.aspx?category=120",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "http://www.laprensa.com.ar/category.aspx?category=150",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicScience,
			contenttopics.ContentTopicHealth,
		},
	}, {
		URL: "http://www.laprensa.com.ar/category.aspx?category=160",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicLiterature,
		},
	}, {
		URL: "https://www.emol.com/nacional",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsChile,
		},
	}, {
		URL: "https://www.emol.com/internacional",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.emol.com/tecnologia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "https://www.emol.com/economia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://www.emol.com/deportes",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://www.emol.com/espectaculos",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://www.emol.com/tendencias/salud",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicHealth,
		},
	}, {
		URL: "https://www.emol.com/tendencias/belleza",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFashion,
		},
	}, {
		URL: "https://www.emol.com/autos",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicAutomotive,
		},
	}, {
		URL: "https://www.latercera.com/canal/el-deportivo/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://www.latercera.com/canal/politica",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicPolitics,
		},
	}, {
		URL: "https://www.latercera.com/etiqueta/empresas-mercados",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicBusiness,
		},
	}, {
		URL: "https://www.latercera.com/etiqueta/economia-dinero",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://www.latercera.com/canal/pulso",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
			contenttopics.ContentTopicBusiness,
		},
	}, {
		URL: "https://www.latercera.com/etiqueta/practico-tecnologia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "https://www.latercera.com/etiqueta/practico-belleza-y-salud/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFashion,
		},
	}, {
		URL: "https://www.latercera.com/etiqueta/practico-casa-y-cocina/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicHome,
			contenttopics.ContentTopicCooking,
		},
	}, {
		URL: "https://www.latercera.com/canal/nacional",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsChile,
		},
	}, {
		URL: "https://www.latercera.com/canal/mundo",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.latercera.com/canal/opinion",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "http://glamorama.latercera.com/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCelebrityNews,
		},
	}, {
		URL: "https://www.latercera.com/etiqueta/libros",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicLiterature,
		},
	}, {
		URL: "https://www.latercera.com/etiqueta/musica",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicMusic,
		},
	}, {
		URL: "https://www.latercera.com/etiqueta/cine/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFilm,
		},
	}, {
		URL: "https://www.latercera.com/etiqueta/videojuegos-de-culto",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicVideoGames,
		},
	}, {
		URL: "https://www.latercera.com/etiqueta/arte",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArt,
		},
	},
}
