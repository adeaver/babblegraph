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
	}, {
		URL: "https://www.lacuarta.com/canal/mundo",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.lacuarta.com/canal/deportes",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://www.lacuarta.com/canal/espectaculos",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://www.elcolombiano.com/colombia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsColombia,
		},
	}, {
		URL: "https://www.elcolombiano.com/internacional",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.elcolombiano.com/negocios",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicBusiness,
		},
	}, {
		URL: "https://www.elcolombiano.com/negocios/economia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://www.elcolombiano.com/negocios/finanzas",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFinance,
		},
	}, {
		URL: "https://www.elcolombiano.com/deportes",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://www.elcolombiano.com/opinion",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://www.elcolombiano.com/cultura/cine",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFilm,
		},
	}, {
		URL: "https://www.elcolombiano.com/cultura/literatura",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicLiterature,
		},
	}, {
		URL: "https://www.elcolombiano.com/cultura/musica",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicMusic,
		},
	}, {
		URL: "https://www.elcolombiano.com/tecnologia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "https://www.elcolombiano.com/tecnologia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "https://www.elcolombiano.com/tecnologia/videojuegos",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicVideoGames,
		},
	}, {
		URL: "https://www.elcolombiano.com/entretenimiento",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://www.elbogotano.com.co/category/poder",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicPolitics,
		},
	}, {
		URL: "https://www.elbogotano.com.co/category/noticias",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsColombia,
		},
	}, {
		URL: "https://www.elespectador.com/opinion",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://www.elespectador.com/noticias/economia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://www.elespectador.com/noticias/tecnologia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "https://www.elespectador.com/entretenimiento",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://www.elespectador.com/entretenimiento/cine",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFilm,
		},
	}, {
		URL: "https://www.elespectador.com/entretenimiento/musica",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicMusic,
		},
	}, {
		URL: "https://www.elespectador.com/deportes",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://www.elespectador.com/noticias/autos",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicAutomotive,
		},
	}, {
		URL: "https://www.elespectador.com/noticias/politica",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicPolitics,
		},
	}, {
		URL: "https://www.elespectador.com/noticias/nacional",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsColombia,
		},
	}, {
		URL: "https://www.elespectador.com/noticias/el-mundo",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.elespectador.com/noticias/ciencia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicScience,
		},
	}, {
		URL: "https://www.elespectador.com/noticias/salud",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicHealth,
		},
	}, {
		URL: "https://www.nacion.com/el-pais",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsCostaRica,
		},
	}, {
		URL: "https://www.nacion.com/puro-deporte",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://www.nacion.com/economia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://www.nacion.com/opinion",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://www.nacion.com/viva/cine",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFilm,
		},
	}, {
		URL: "https://www.nacion.com/viva/entretenimiento",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://www.nacion.com/viva/musica",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicMusic,
		},
	}, {
		URL: "https://www.nacion.com/viva/moda",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFashion,
		},
	}, {
		URL: "https://www.nacion.com/el-mundo",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.nacion.com/ciencia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicScience,
		},
	}, {
		URL: "https://www.nacion.com/ciencia/salud",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicHealth,
		},
	}, {
		URL: "https://www.nacion.com/tecnologia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "https://www.nacion.com/tecnologia/videojuegos",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicVideoGames,
		},
	}, {
		URL: "https://www.nacion.com/sabores",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCooking,
		},
	}, {
		URL: "https://www.elsalvador.com/category/noticias/internacional",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.elsalvador.com/category/noticias/nacional",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsElSalvador,
		},
	}, {
		URL: "https://www.elsalvador.com/category/entretenimiento/espectaculos",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://www.elsalvador.com/category/entretenimiento/tecnologia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "https://www.elsalvador.com/category/entretenimiento/turismo",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTravel,
		},
	}, {
		URL: "https://www.elsalvador.com/category/entretenimiento/cultura",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCulture,
		},
	}, {
		URL: "https://www.elsalvador.com/category/vida/salud",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicHealth,
		},
	}, {
		URL: "https://www.elsalvador.com/category/deportes",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://www.elsalvador.com/category/opinion",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://diario.elmundo.sv/category/politica",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicPolitics,
		},
	}, {
		URL: "https://diario.elmundo.sv/category/economia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://diario.elmundo.sv/category/nacionales",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsElSalvador,
		},
	}, {
		URL: "https://diario.elmundo.sv/category/internacionales",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://diario.elmundo.sv/category/entretenimiento",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://diario.elmundo.sv/category/deportes",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://diario.elmundo.sv/category/empresarial",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicBusiness,
		},
	}, {
		URL: "https://diario.elmundo.sv/category/opinion",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://dca.gob.gt/noticias-guatemala-diario-centro-america/category/nacionales",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsGuatemala,
		},
	}, {
		URL: "https://dca.gob.gt/noticias-guatemala-diario-centro-america/category/opinion",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://dca.gob.gt/noticias-guatemala-diario-centro-america/category/deportes",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://dca.gob.gt/noticias-guatemala-diario-centro-america/category/artes",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCulture,
		},
	}, {
		URL: "https://dca.gob.gt/noticias-guatemala-diario-centro-america/category/economicas",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://dca.gob.gt/noticias-guatemala-diario-centro-america/category/internacionales/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://dca.gob.gt/noticias-guatemala-diario-centro-america/category/salud/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicHealth,
		},
	}, {
		URL: "https://www.prensalibre.com/guatemala",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsGuatemala,
		},
	}, {
		URL: "https://www.prensalibre.com/deportes",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://www.prensalibre.com/internacional",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.prensalibre.com/economia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://www.prensalibre.com/vida",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicLifestyle,
		},
	}, {
		URL: "https://www.prensalibre.com/opinion",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://www.prensalibre.com/vida/tecnologia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "https://www.prensalibre.com/vida/salud-y-familia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicHealth,
		},
	}, {
		URL: "https://www.laprensa.hn/honduras",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsHonduras,
		},
	}, {
		URL: "http://www.laprensa.hn/deportes",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "http://www.laprensa.hn/espectaculos",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://www.laprensa.hn/mundo",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "http://www.laprensa.hn/tecnologia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "http://www.laprensa.hn/opinion",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://jornada.com.mx/category/politica.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicPolitics,
		},
	}, {
		URL: "https://jornada.com.mx/category/economia.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://jornada.com.mx/category/mundo.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://jornada.com.mx/category/estados.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsMexico,
		},
	}, {
		URL: "https://jornada.com.mx/category/deportes.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://jornada.com.mx/category/cultura.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://jornada.com.mx/category/ciencia-y-tecnologia.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicScience,
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "https://jornada.com.mx/category/autos.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicAutomotive,
		},
	}, {
		URL: "https://heraldodemexico.com.mx/nacional",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsMexico,
		},
	}, {
		URL: "https://heraldodemexico.com.mx/mundo/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://heraldodemexico.com.mx/economia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://heraldodemexico.com.mx/espectaculos/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://heraldodemexico.com.mx/temas/viajes-1865.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTravel,
		},
	}, {
		URL: "https://heraldodemexico.com.mx/temas/autos-1868.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicAutomotive,
		},
	}, {
		URL: "https://heraldodemexico.com.mx/cultura/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCulture,
		},
	}, {
		URL: "https://heraldodemexico.com.mx/tecnologia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "https://www.eleconomista.com.mx/seccion/empresas/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicBusiness,
		},
	}, {
		URL: "https://www.eleconomista.com.mx/seccion/sector-financiero/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFinance,
		},
	}, {
		URL: "https://www.eleconomista.com.mx/seccion/internacionales/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.eleconomista.com.mx/seccion/tecnologia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "https://www.eleconomista.com.mx/seccion/opinion/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://www.eleconomista.com.mx/seccion/arte-ideas/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCulture,
		},
	}, {
		URL: "https://www.eleconomista.com.mx/seccion/economia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://www.eleconomista.com.mx/seccion/politica/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicPolitics,
		},
	}, {
		URL: "https://www.eleconomista.com.mx/seccion/deportes/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://www.laprensa.com.ni/nacionales",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsNicaragua,
		},
	}, {
		URL: "https://www.laprensa.com.ni/politica",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicPolitics,
		},
	}, {
		URL: "https://www.laprensa.com.ni/economia",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://www.laprensa.com.ni/deportes",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://www.laprensa.com.ni/internacionales",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.laprensa.com.ni/espectaculo",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://www.laprensa.com.ni/suplemento/empresariales",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicBusiness,
		},
	}, {
		URL: "https://www.laprensa.com.ni/opinion-main",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://www.laprensa.com.ni/salud",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicHealth,
		},
	}, {
		URL: "https://www.lajornadanet.com/index.php/noticias/nicaragua",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsNicaragua,
		},
	}, {
		URL: "https://www.lajornadanet.com/index.php/noticias/nicaragua/politica",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicPolitics,
		},
	}, {
		URL: "https://www.lajornadanet.com/index.php/noticias/empresariales/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicBusiness,
		},
	}, {
		URL: "https://www.lajornadanet.com/index.php/noticias/economia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://www.lajornadanet.com/index.php/noticias/mundo/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.lajornadanet.com/index.php/noticias/ciencia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicScience,
		},
	}, {
		URL: "https://www.lajornadanet.com/index.php/noticias/mexico/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsMexico,
		},
	}, {
		URL: "https://www.lajornadanet.com/index.php/noticias/costa-rica/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsCostaRica,
		},
	}, {
		URL: "https://www.lajornadanet.com/index.php/noticias/tecno/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "http://elsiglo.com.pa/tag/panama/13",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsPanama,
		},
	}, {
		URL: "http://elsiglo.com.pa/tag/internacional/17",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "http://elsiglo.com.pa/tag/economia/18",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "http://elsiglo.com.pa/tag/deportes/19",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "http://elsiglo.com.pa/tag/espectaculos/20",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "http://elsiglo.com.pa/tag/opinion/22",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://larepublica.pe/politica/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicPolitics,
		},
	}, {
		URL: "https://larepublica.pe/economia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://larepublica.pe/mundo/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://larepublica.pe/deportes/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://larepublica.pe/salud/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicHealth,
		},
	}, {
		URL: "https://larepublica.pe/cultural/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCulture,
		},
	}, {
		URL: "https://larepublica.pe/ciencia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicScience,
		},
	}, {
		URL: "https://larepublica.pe/turismo/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTravel,
		},
	}, {
		URL: "https://larepublica.pe/espectaculos/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://larepublica.pe/tecnologia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "https://larepublica.pe/cine-series/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFilm,
		},
	}, {
		URL: "https://larepublica.pe/videojuegos/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicVideoGames,
		},
	}, {
		URL: "https://larepublica.pe/estilo/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFashion,
		},
	}, {
		URL: "https://larepublica.pe/region-norte/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsPeru,
		},
	}, {
		URL: "https://larepublica.pe/region-sur/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsPeru,
		},
	}, {
		URL: "https://www.abc.com.py/nacionales",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsParaguay,
		},
	}, {
		URL: "https://www.abc.com.py/deportes",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://www.abc.com.py/espectaculos",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://www.abc.com.py/internacionales",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.ultimahora.com/contenidos/nacional.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsParaguay,
		},
	}, {
		URL: "https://d10.ultimahora.com/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://tvo.ultimahora.com/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCelebrityNews,
		},
	}, {
		URL: "https://www.ultimahora.com/contenidos/gaming.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicVideoGames,
		},
	}, {
		URL: "https://www.ultimahora.com/contenidos/arte-y-espectaculos.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEntertainment,
		},
	}, {
		URL: "https://www.ultimahora.com/contenidos/mundo.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.ultimahora.com/contenidos/turismo.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTravel,
		},
	}, {
		URL: "https://www.abc.es/espana/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsSpain,
		},
	}, {
		URL: "https://www.abc.es/internacional",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.abc.es/economia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://www.abc.es/opinion/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://www.abc.es/deportes/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://www.abc.es/estilo/moda/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFashion,
		},
	}, {
		URL: "https://www.abc.es/cultura/libros/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicLiterature,
		},
	}, {
		URL: "https://www.abc.es/cultura/musica/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicMusic,
		},
	}, {
		URL: "https://www.abc.es/cultura/arte/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArt,
		},
	}, {
		URL: "https://www.abc.es/cultura/teatros/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTheater,
		},
	}, {
		URL: "https://www.abc.es/ciencia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicScience,
		},
	}, {
		URL: "https://www.abc.es/viajar/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTravel,
		},
	}, {
		URL: "https://www.elmundo.es/espana.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsSpain,
		},
	}, {
		URL: "https://www.elmundo.es/opinion.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://www.elmundo.es/economia.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicEconomy,
		},
	}, {
		URL: "https://www.elmundo.es/internacional.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://www.elmundo.es/deportes.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://www.elmundo.es/cultura/cine.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFilm,
		},
	}, {
		URL: "https://www.elmundo.es/cultura/literatura.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicLiterature,
		},
	}, {
		URL: "https://www.elmundo.es/cultura/musica.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicMusic,
		},
	}, {
		URL: "https://www.elmundo.es/cultura/teatro.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTheater,
		},
	}, {
		URL: "https://www.elmundo.es/cultura/arte.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArt,
		},
	}, {
		URL: "https://www.elmundo.es/ciencia-y-salud/ciencia.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicScience,
		},
	}, {
		URL: "https://www.elmundo.es/ciencia-y-salud/salud.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicHealth,
		},
	}, {
		URL: "https://www.elmundo.es/tecnologia.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "https://www.elmundo.es/tecnologia/videojuegos.html",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicVideoGames,
		},
	}, {
		URL: "https://elpais.com/internacional/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicWorldNews,
		},
	}, {
		URL: "https://elpais.com/opinion/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicOpinion,
		},
	}, {
		URL: "https://elpais.com/espana/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicCurrentEventsSpain,
		},
	}, {
		URL: "https://elpais.com/ciencia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicScience,
		},
	}, {
		URL: "https://elpais.com/tecnologia/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTechnology,
		},
	}, {
		URL: "https://elpais.com/noticias/libros/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicLiterature,
		},
	}, {
		URL: "https://elpais.com/noticias/cine/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicFilm,
		},
	}, {
		URL: "https://elpais.com/noticias/musica/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicMusic,
		},
	}, {
		URL: "https://elpais.com/noticias/teatro/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicTheater,
		},
	}, {
		URL: "https://elpais.com/noticias/arte/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArt,
		},
	}, {
		URL: "https://elpais.com/noticias/arquitectura/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicArchitecture,
		},
	}, {
		URL: "https://elpais.com/deportes/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicSports,
		},
	}, {
		URL: "https://motor.elpais.com/",
		Topics: []contenttopics.ContentTopic{
			contenttopics.ContentTopicAutomotive,
		},
	},
}
