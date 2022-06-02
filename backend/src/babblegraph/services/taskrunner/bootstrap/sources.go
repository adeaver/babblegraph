package bootstrap

import (
	"babblegraph/util/geo"
	"babblegraph/wordsmith"
)

type allowableDomain struct {
	Country      geo.CountryCode
	LanguageCode wordsmith.LanguageCode
	SeedURLs     []seedURL
}

type seedURL struct {
	URL        string
	TopicLabel string
}

var Sources = map[string]allowableDomain{
	"ambito.com": {
		Country:      geo.CountryCodeArgentina,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.ambito.com/contenidos/economia.html",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://www.ambito.com/contenidos/finanzas.html",
				TopicLabel: "finance",
			},
			seedURL{
				URL:        "https://www.ambito.com/contenidos/politica.html",
				TopicLabel: "politics",
			},
			seedURL{
				URL:        "https://www.ambito.com/contenidos/negocios.html",
				TopicLabel: "business",
			},
			seedURL{
				URL:        "https://www.ambito.com/contenidos/lifestyle.html",
				TopicLabel: "lifestyle",
			},
			seedURL{
				URL:        "https://www.ambito.com/contenidos/opinion.html",
				TopicLabel: "opinion",
			},
			seedURL{
				URL:        "https://www.ambito.com/contenidos/nacional.html",
				TopicLabel: "current-events-argentina",
			},
			seedURL{
				URL:        "https://www.ambito.com/contenidos/mundo.html",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.ambito.com/contenidos/espectaculos.html",
				TopicLabel: "celebrity-news",
			},
			seedURL{
				URL:        "https://www.ambito.com/contenidos/autos.html",
				TopicLabel: "automotive",
			},
			seedURL{
				URL:        "https://www.ambito.com/contenidos/deportes.html",
				TopicLabel: "sports",
			},
		},
	},
	"cronista.com": allowableDomain{
		Country:      geo.CountryCodeArgentina,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.cronista.com/apertura-negocio",
				TopicLabel: "business",
			},
			seedURL{
				URL:        "https://www.cronista.com/seccion/economia_politica",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://www.cronista.com/seccion/finanzas_mercados",
				TopicLabel: "finance",
			},
			seedURL{
				URL:        "https://www.cronista.com/seccion/internacionales",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.cronista.com/seccion/columnistas",
				TopicLabel: "opinion",
			},
		},
	},
	"emol.com": allowableDomain{
		Country:      geo.CountryCodeChile,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.emol.com/nacional",
				TopicLabel: "current-events-chile",
			},
			seedURL{
				URL:        "https://www.emol.com/internacional",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.emol.com/tecnologia",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://www.emol.com/economia",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://www.emol.com/deportes",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.emol.com/espectaculos",
				TopicLabel: "celebrity-news",
			},
			seedURL{
				URL:        "https://www.emol.com/tendencias/salud",
				TopicLabel: "health",
			},
			seedURL{
				URL:        "https://www.emol.com/tendencias/belleza",
				TopicLabel: "fashion",
			},
			seedURL{
				URL:        "https://www.emol.com/autos",
				TopicLabel: "automotive",
			},
		},
	},
	"latercera.com": allowableDomain{
		Country:      geo.CountryCodeChile,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.latercera.com/canal/el-deportivo/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.latercera.com/canal/politica",
				TopicLabel: "politics",
			},
			seedURL{
				URL:        "https://www.latercera.com/etiqueta/empresas-mercados",
				TopicLabel: "business",
			},
			seedURL{
				URL:        "https://www.latercera.com/etiqueta/economia-dinero",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://www.latercera.com/canal/pulso",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://www.latercera.com/etiqueta/practico-tecnologia/",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://www.latercera.com/etiqueta/practico-belleza-y-salud/",
				TopicLabel: "fashion",
			},
			seedURL{
				URL:        "https://www.latercera.com/etiqueta/practico-casa-y-cocina/",
				TopicLabel: "kitchen",
			},
			seedURL{
				URL:        "https://www.latercera.com/canal/nacional",
				TopicLabel: "current-events-chile",
			},
			seedURL{
				URL:        "https://www.latercera.com/canal/mundo",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.latercera.com/canal/opinion",
				TopicLabel: "opinion",
			},
			seedURL{
				URL:        "http://glamorama.latercera.com/",
				TopicLabel: "fashion",
			},
			seedURL{
				URL:        "https://www.latercera.com/etiqueta/libros",
				TopicLabel: "literature",
			},
			seedURL{
				URL:        "https://www.latercera.com/etiqueta/musica",
				TopicLabel: "music",
			},
			seedURL{
				URL:        "https://www.latercera.com/etiqueta/cine/",
				TopicLabel: "film",
			},
			seedURL{
				URL:        "https://www.latercera.com/etiqueta/videojuegos-de-culto",
				TopicLabel: "video-games",
			},
			seedURL{
				URL:        "https://www.latercera.com/etiqueta/arte",
				TopicLabel: "art",
			},
		},
	},
	"lacuarta.cl": allowableDomain{
		Country:      geo.CountryCodeChile,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	},
	"elbogotano.com.co": allowableDomain{
		Country:      geo.CountryCodeColombia,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.elbogotano.com.co/category/poder",
				TopicLabel: "current-events-colombia",
			},
			seedURL{
				URL:        "https://www.elbogotano.com.co/category/noticias",
				TopicLabel: "current-events-colombia",
			},
		},
	},
	"elespectador.com": allowableDomain{
		Country:      geo.CountryCodeColombia,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.elespectador.com/opinion",
				TopicLabel: "opinion",
			},
			seedURL{
				URL:        "https://www.elespectador.com/noticias/economia",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://www.elespectador.com/noticias/tecnologia",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://www.elespectador.com/entretenimiento",
				TopicLabel: "entertainment",
			},
			seedURL{
				URL:        "https://www.elespectador.com/entretenimiento/cine",
				TopicLabel: "film",
			},
			seedURL{
				URL:        "https://www.elespectador.com/entretenimiento/musica",
				TopicLabel: "music",
			},
			seedURL{
				URL:        "https://www.elespectador.com/deportes",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.elespectador.com/noticias/autos",
				TopicLabel: "automotive",
			},
			seedURL{
				URL:        "https://www.elespectador.com/noticias/politica",
				TopicLabel: "politics",
			},
			seedURL{
				URL:        "https://www.elespectador.com/noticias/nacional",
				TopicLabel: "current-events-colombia",
			},
			seedURL{
				URL:        "https://www.elespectador.com/noticias/el-mundo",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.elespectador.com/noticias/ciencia",
				TopicLabel: "science",
			},
			seedURL{
				URL:        "https://www.elespectador.com/noticias/salud",
				TopicLabel: "health",
			},
		},
	},
	"elsalvador.com": allowableDomain{
		Country:      geo.CountryCodeElSalvador,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.elsalvador.com/category/noticias/internacional",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.elsalvador.com/category/noticias/nacional",
				TopicLabel: "current-events-el-salvador",
			},
			seedURL{
				URL:        "https://www.elsalvador.com/category/entretenimiento/espectaculos",
				TopicLabel: "celebrity-news",
			},
			seedURL{
				URL:        "https://www.elsalvador.com/category/entretenimiento/tecnologia",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://www.elsalvador.com/category/entretenimiento/turismo",
				TopicLabel: "travel",
			},
			seedURL{
				URL:        "https://www.elsalvador.com/category/entretenimiento/cultura",
				TopicLabel: "culture",
			},
			seedURL{
				URL:        "https://www.elsalvador.com/category/vida/salud",
				TopicLabel: "health",
			},
			seedURL{
				URL:        "https://www.elsalvador.com/category/deportes",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.elsalvador.com/category/opinion",
				TopicLabel: "opinion",
			},
		},
	},
	"elmundo.sv": allowableDomain{
		Country:      geo.CountryCodeElSalvador,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	},
	"dca.gob.gt": allowableDomain{
		Country:      geo.CountryCodeGuatemala,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://dca.gob.gt/noticias-guatemala-diario-centro-america/category/nacionales",
				TopicLabel: "current-events-guatemala",
			},
			seedURL{
				URL:        "https://dca.gob.gt/noticias-guatemala-diario-centro-america/category/opinion",
				TopicLabel: "opinion",
			},
			seedURL{
				URL:        "https://dca.gob.gt/noticias-guatemala-diario-centro-america/category/deportes",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://dca.gob.gt/noticias-guatemala-diario-centro-america/category/artes",
				TopicLabel: "art",
			},
			seedURL{
				URL:        "https://dca.gob.gt/noticias-guatemala-diario-centro-america/category/economicas",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://dca.gob.gt/noticias-guatemala-diario-centro-america/category/internacionales/",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://dca.gob.gt/noticias-guatemala-diario-centro-america/category/salud/",
				TopicLabel: "health",
			},
		},
	},
	"prensalibre.com": allowableDomain{
		Country:      geo.CountryCodeGuatemala,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.prensalibre.com/guatemala",
				TopicLabel: "current-events-guatemala",
			},
			seedURL{
				URL:        "https://www.prensalibre.com/deportes",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.prensalibre.com/internacional",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.prensalibre.com/economia",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://www.prensalibre.com/vida",
				TopicLabel: "lifestyle",
			},
			seedURL{
				URL:        "https://www.prensalibre.com/opinion",
				TopicLabel: "opinion",
			},
			seedURL{
				URL:        "https://www.prensalibre.com/vida/tecnologia/",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://www.prensalibre.com/vida/salud-y-familia/",
				TopicLabel: "health",
			},
		},
	},
	"heraldodemexico.com.mx": allowableDomain{
		Country:      geo.CountryCodeMexico,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://heraldodemexico.com.mx/nacional",
				TopicLabel: "current-events-mexico",
			},
			seedURL{
				URL:        "https://heraldodemexico.com.mx/mundo/",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://heraldodemexico.com.mx/economia/",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://heraldodemexico.com.mx/espectaculos/",
				TopicLabel: "entertainment",
			},
			seedURL{
				URL:        "https://heraldodemexico.com.mx/temas/viajes-1865.html",
				TopicLabel: "travel",
			},
			seedURL{
				URL:        "https://heraldodemexico.com.mx/temas/autos-1868.html",
				TopicLabel: "automotive",
			},
			seedURL{
				URL:        "https://heraldodemexico.com.mx/cultura/",
				TopicLabel: "culture",
			},
			seedURL{
				URL:        "https://heraldodemexico.com.mx/tecnologia/",
				TopicLabel: "technology",
			},
		},
	},
	"eleconomista.com.mx": allowableDomain{
		Country:      geo.CountryCodeMexico,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.eleconomista.com.mx/seccion/empresas/",
				TopicLabel: "business",
			},
			seedURL{
				URL:        "https://www.eleconomista.com.mx/seccion/sector-financiero/",
				TopicLabel: "finance",
			},
			seedURL{
				URL:        "https://www.eleconomista.com.mx/seccion/internacionales/",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.eleconomista.com.mx/seccion/tecnologia/",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://www.eleconomista.com.mx/seccion/opinion/",
				TopicLabel: "opinion",
			},
			seedURL{
				URL:        "https://www.eleconomista.com.mx/seccion/arte-ideas/",
				TopicLabel: "art",
			},
			seedURL{
				URL:        "https://www.eleconomista.com.mx/seccion/economia/",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://www.eleconomista.com.mx/seccion/politica/",
				TopicLabel: "politics",
			},
			seedURL{
				URL:        "https://www.eleconomista.com.mx/seccion/deportes/",
				TopicLabel: "sports",
			},
		},
	},
	"yucatan.com.mx": allowableDomain{
		Country:      geo.CountryCodeMexico,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.yucatan.com.mx/seccion/mexico",
				TopicLabel: "current-events-mexico",
			},
			seedURL{
				URL:        "https://www.yucatan.com.mx/seccion/internacional",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.yucatan.com.mx/seccion/deportes",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.yucatan.com.mx/seccion/espectaculos",
				TopicLabel: "celebrity-news",
			},
			seedURL{
				URL:        "https://www.yucatan.com.mx/seccion/merida",
				TopicLabel: "current-events-mexico",
			},
			seedURL{
				URL:        "https://www.yucatan.com.mx/seccion/yucatan",
				TopicLabel: "current-events-mexico",
			},
		},
	},
	"informador.mx": allowableDomain{
		Country:      geo.CountryCodeMexico,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.informador.mx/seccion/jalisco/",
				TopicLabel: "current-events-mexico",
			},
			seedURL{
				URL:        "https://www.informador.mx/seccion/mexico/",
				TopicLabel: "current-events-mexico",
			},
			seedURL{
				URL:        "https://www.informador.mx/seccion/internacional/",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.informador.mx/seccion/deportes/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.informador.mx/seccion/entretenimiento/",
				TopicLabel: "entertainment",
			},
			seedURL{
				URL:        "https://www.informador.mx/autos-t2294",
				TopicLabel: "automotive",
			},
			seedURL{
				URL:        "https://www.informador.mx/seccion/economia/",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://www.informador.mx/seccion/tecnologia/",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://www.informador.mx/seccion/cultura/",
				TopicLabel: "culture",
			},
		},
	},
	"laprensa.com.ni": allowableDomain{
		Country:      geo.CountryCodeNicaragua,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	},
	"lajornadanet.com": allowableDomain{
		Country:      geo.CountryCodeNicaragua,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.lajornadanet.com/index.php/noticias/nicaragua",
				TopicLabel: "current-events-nicaragua",
			},
			seedURL{
				URL:        "https://www.lajornadanet.com/index.php/noticias/nicaragua/politica",
				TopicLabel: "politics",
			},
			seedURL{
				URL:        "https://www.lajornadanet.com/index.php/noticias/empresariales/",
				TopicLabel: "business",
			},
			seedURL{
				URL:        "https://www.lajornadanet.com/index.php/noticias/economia/",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://www.lajornadanet.com/index.php/noticias/mundo/",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.lajornadanet.com/index.php/noticias/ciencia/",
				TopicLabel: "science",
			},
			seedURL{
				URL:        "https://www.lajornadanet.com/index.php/noticias/mexico/",
				TopicLabel: "current-events-mexico",
			},
			seedURL{
				URL:        "https://www.lajornadanet.com/index.php/noticias/costa-rica/",
				TopicLabel: "current-events-costa-rica",
			},
			seedURL{
				URL:        "https://www.lajornadanet.com/index.php/noticias/tecno/",
				TopicLabel: "technology",
			},
		},
	},
	"elsiglo.com": allowableDomain{
		Country:      geo.CountryCodePanama,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "http://elsiglo.com.pa/tag/panama/13",
				TopicLabel: "current-events-panama",
			},
			seedURL{
				URL:        "http://elsiglo.com.pa/tag/internacional/17",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "http://elsiglo.com.pa/tag/economia/18",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "http://elsiglo.com.pa/tag/deportes/19",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "http://elsiglo.com.pa/tag/espectaculos/20",
				TopicLabel: "entertainment",
			},
			seedURL{
				URL:        "http://elsiglo.com.pa/tag/opinion/22",
				TopicLabel: "opinion",
			},
		},
	},
	"larepublica.pe": allowableDomain{
		Country:      geo.CountryCodePeru,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://larepublica.pe/politica/",
				TopicLabel: "politics",
			},
			seedURL{
				URL:        "https://larepublica.pe/economia/",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://larepublica.pe/mundo/",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://larepublica.pe/deportes/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://larepublica.pe/salud/",
				TopicLabel: "health",
			},
			seedURL{
				URL:        "https://larepublica.pe/cultural/",
				TopicLabel: "culture",
			},
			seedURL{
				URL:        "https://larepublica.pe/ciencia/",
				TopicLabel: "science",
			},
			seedURL{
				URL:        "https://larepublica.pe/turismo/",
				TopicLabel: "travel",
			},
			seedURL{
				URL:        "https://larepublica.pe/espectaculos/",
				TopicLabel: "enterainment",
			},
			seedURL{
				URL:        "https://larepublica.pe/tecnologia/",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://larepublica.pe/cine-series/",
				TopicLabel: "film",
			},
			seedURL{
				URL:        "https://larepublica.pe/videojuegos/",
				TopicLabel: "video-games",
			},
			seedURL{
				URL:        "https://larepublica.pe/estilo/",
				TopicLabel: "fashion",
			},
			seedURL{
				URL:        "https://larepublica.pe/region-norte/",
				TopicLabel: "current-events-peru",
			},
			seedURL{
				URL:        "https://larepublica.pe/region-sur/",
				TopicLabel: "current-events-peru",
			},
		},
	},
	"abc.com.py": allowableDomain{
		Country:      geo.CountryCodeParaguay,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	},
	"ultimahora.com": allowableDomain{
		Country:      geo.CountryCodeParaguay,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.ultimahora.com/contenidos/nacional.html",
				TopicLabel: "current-events-paraguay",
			},
			seedURL{
				URL:        "https://d10.ultimahora.com/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://tvo.ultimahora.com/",
				TopicLabel: "enterainment",
			},
			seedURL{
				URL:        "https://www.ultimahora.com/contenidos/gaming.html",
				TopicLabel: "video-games",
			},
			seedURL{
				URL:        "https://www.ultimahora.com/contenidos/arte-y-espectaculos.html",
				TopicLabel: "art",
			},
			seedURL{
				URL:        "https://www.ultimahora.com/contenidos/mundo.html",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.ultimahora.com/contenidos/turismo.html",
				TopicLabel: "travel",
			},
		},
	},
	"abc.es": allowableDomain{
		Country:      geo.CountryCodeSpain,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.abc.com.py/nacionales",
				TopicLabel: "current-events-paraguay",
			},
			seedURL{
				URL:        "https://www.abc.com.py/deportes",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.abc.com.py/espectaculos",
				TopicLabel: "entertainment",
			},
			seedURL{
				URL:        "https://www.abc.com.py/internacionales",
				TopicLabel: "world-news",
			},
		},
	},
	"elmundo.es": allowableDomain{
		Country:      geo.CountryCodeSpain,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.elmundo.es/espana.html",
				TopicLabel: "current-events-spain",
			},
			seedURL{
				URL:        "https://www.elmundo.es/opinion.html",
				TopicLabel: "opinion",
			},
			seedURL{
				URL:        "https://www.elmundo.es/economia.html",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://www.elmundo.es/internacional.html",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.elmundo.es/deportes.html",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.elmundo.es/cultura/cine.html",
				TopicLabel: "film",
			},
			seedURL{
				URL:        "https://www.elmundo.es/cultura/literatura.html",
				TopicLabel: "literature",
			},
			seedURL{
				URL:        "https://www.elmundo.es/cultura/musica.html",
				TopicLabel: "music",
			},
			seedURL{
				URL:        "https://www.elmundo.es/cultura/teatro.html",
				TopicLabel: "theater",
			},
			seedURL{
				URL:        "https://www.elmundo.es/cultura/arte.html",
				TopicLabel: "art",
			},
			seedURL{
				URL:        "https://www.elmundo.es/ciencia-y-salud/ciencia.html",
				TopicLabel: "science",
			},
			seedURL{
				URL:        "https://www.elmundo.es/ciencia-y-salud/salud.html",
				TopicLabel: "health",
			},
			seedURL{
				URL:        "https://www.elmundo.es/tecnologia.html",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://www.elmundo.es/tecnologia/videojuegos.html",
				TopicLabel: "video-games",
			},
		},
	},
	"elpais.com": allowableDomain{
		Country:      geo.CountryCodeSpain,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://elpais.com/internacional/",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://elpais.com/opinion/",
				TopicLabel: "opinion",
			},
			seedURL{
				URL:        "https://elpais.com/espana/",
				TopicLabel: "current-events-spain",
			},
			seedURL{
				URL:        "https://elpais.com/ciencia/",
				TopicLabel: "science",
			},
			seedURL{
				URL:        "https://elpais.com/tecnologia/",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://elpais.com/noticias/libros/",
				TopicLabel: "literature",
			},
			seedURL{
				URL:        "https://elpais.com/noticias/cine/",
				TopicLabel: "film",
			},
			seedURL{
				URL:        "https://elpais.com/noticias/musica/",
				TopicLabel: "music",
			},
			seedURL{
				URL:        "https://elpais.com/noticias/teatro/",
				TopicLabel: "theater",
			},
			seedURL{
				URL:        "https://elpais.com/noticias/arte/",
				TopicLabel: "art",
			},
			seedURL{
				URL:        "https://elpais.com/noticias/arquitectura/",
				TopicLabel: "architecture",
			},
			seedURL{
				URL:        "https://elpais.com/deportes/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://motor.elpais.com/",
				TopicLabel: "automotive",
			},
			seedURL{
				URL:        "https://elviajero.elpais.com/",
				TopicLabel: "travel",
			},
			seedURL{
				URL:        "https://elpais.com/economia/negocios/",
				TopicLabel: "business",
			},
		},
	},
	"elpais.com.uy": allowableDomain{
		Country:      geo.CountryCodeUruguay,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.elpais.com.uy/mundo",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.elpais.com.uy/opinion",
				TopicLabel: "opinion",
			},
			seedURL{
				URL:        "https://www.elpais.com.uy/negocios",
				TopicLabel: "business",
			},
			seedURL{
				URL:        "https://www.elpais.com.uy/ovacion",
				TopicLabel: "opinion",
			},
		},
	},
	"elcomercio.es": allowableDomain{
		Country:      geo.CountryCodeSpain,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.elcomercio.es/asturias/",
				TopicLabel: "current-events-spain",
			},
			seedURL{
				URL:        "https://www.elcomercio.es/politica/",
				TopicLabel: "politics",
			},
			seedURL{
				URL:        "https://www.elcomercio.es/internacional/",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.elcomercio.es/economia/",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://www.elcomercio.es/deportes/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.elcomercio.es/culturas/cine/",
				TopicLabel: "film",
			},
			seedURL{
				URL:        "https://www.elcomercio.es/culturas/libros/",
				TopicLabel: "literature",
			},
			seedURL{
				URL:        "https://www.elcomercio.es/culturas/arte/",
				TopicLabel: "art",
			},
			seedURL{
				URL:        "https://www.elcomercio.es/culturas/musica/",
				TopicLabel: "music",
			},
		},
	},
	"diariolasamericas.com": allowableDomain{
		Country:      geo.CountryCodeUnitedStates,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.diariolasamericas.com/contenidos/mundo.html",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.diariolasamericas.com/contenidos/eeuu.html",
				TopicLabel: "current-events-united-states",
			},
			seedURL{
				URL:        "https://www.diariolasamericas.com/contenidos/deportes.html",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.diariolasamericas.com/contenidos/cultura.html",
				TopicLabel: "culture",
			},
			seedURL{
				URL:        "https://www.diariolasamericas.com/contenidos/opini%C3%B3n.html",
				TopicLabel: "opinion",
			},
			seedURL{
				URL:        "https://www.diariolasamericas.com/contenidos/tecnologia.html",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://www.diariolasamericas.com/contenidos/economia.html",
				TopicLabel: "economy",
			},
		},
	},
	"laopinion.com": allowableDomain{
		Country:      geo.CountryCodeUnitedStates,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://laopinion.com/categoria/entretenimiento/",
				TopicLabel: "entertainment",
			},
			seedURL{
				URL:        "https://laopinion.com/categoria/deportes/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://laopinion.com/categoria/autos/",
				TopicLabel: "automotive",
			},
			seedURL{
				URL:        "https://laopinion.com/categoria/estilo/",
				TopicLabel: "lifestyle",
			},
			seedURL{
				URL:        "https://laopinion.com/categoria/dinero/",
				TopicLabel: "finance",
			},
			seedURL{
				URL:        "https://laopinion.com/categoria/tecnologia/",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://laopinion.com/categoria/opinion/",
				TopicLabel: "opinion",
			},
			seedURL{
				URL:        "https://laopinion.com/categoria-guia-de-compras/salud/",
				TopicLabel: "health",
			},
			seedURL{
				URL:        "https://laopinion.com/categoria-guia-de-compras/tecnologia/",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://laopinion.com/categoria-guia-de-compras/ropa-y-accesorios/",
				TopicLabel: "fashion",
			},
		},
	},
	"eldiariony.com": allowableDomain{
		Country:      geo.CountryCodeUnitedStates,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://eldiariony.com/categoria/entretenimiento/",
				TopicLabel: "entertainment",
			},
			seedURL{
				URL:        "https://eldiariony.com/categoria/deportes/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://eldiariony.com/categoria/salud/",
				TopicLabel: "health",
			},
			seedURL{
				URL:        "https://eldiariony.com/categoria/comida/",
				TopicLabel: "cooking",
			},
			seedURL{
				URL:        "https://eldiariony.com/categoria/estilo/",
				TopicLabel: "lifestyle",
			},
			seedURL{
				URL:        "https://eldiariony.com/categoria/dinero/",
				TopicLabel: "finance",
			},
			seedURL{
				URL:        "https://eldiariony.com/categoria/tecnologia/",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://eldiariony.com/categoria/opinion/",
				TopicLabel: "opinion",
			},
		},
	},
	"hoylosangeles.com": allowableDomain{
		Country:      geo.CountryCodeUnitedStates,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.hoylosangeles.com/espectaculos/musica",
				TopicLabel: "music",
			},
			seedURL{
				URL:        "https://www.hoylosangeles.com/espectaculos/cine",
				TopicLabel: "film",
			},
			seedURL{
				URL:        "https://www.hoylosangeles.com/espectaculos",
				TopicLabel: "entertainment",
			},
			seedURL{
				URL:        "https://www.hoylosangeles.com/deportes",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.hoylosangeles.com/vidayestilo/salud",
				TopicLabel: "health",
			},
			seedURL{
				URL:        "https://www.hoylosangeles.com/vidayestilo/cienciaytecnologia",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://www.hoylosangeles.com/vidayestilo/viajes",
				TopicLabel: "travel",
			},
			seedURL{
				URL:        "https://www.hoylosangeles.com/vidayestilo/autos",
				TopicLabel: "automotive",
			},
			seedURL{
				URL:        "https://www.hoylosangeles.com/vidayestilo/cocina",
				TopicLabel: "cooking",
			},
			seedURL{
				URL:        "https://www.hoylosangeles.com/vidayestilo/libros",
				TopicLabel: "literature",
			},
			seedURL{
				URL:        "https://www.hoylosangeles.com/opinion",
				TopicLabel: "opinion",
			},
			seedURL{
				URL:        "https://www.hoylosangeles.com/noticias/estadosunidos",
				TopicLabel: "current-events-united-states",
			},
		},
	},
	"elnacional.com": allowableDomain{
		Country:      geo.CountryCodeVenezuela,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.elnacional.com/venezuela/",
				TopicLabel: "current-events-venezuela",
			},
			seedURL{
				URL:        "https://www.elnacional.com/mundo/",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.elnacional.com/economia/",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://www.elnacional.com/deportes/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.elnacional.com/salud-ciencia-tecnologia",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://www.elnacional.com/gamers",
				TopicLabel: "video-games",
			},
			seedURL{
				URL:        "https://www.elnacional.com/gadgets",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://www.elnacional.com/cine",
				TopicLabel: "film",
			},
			seedURL{
				URL:        "https://www.elnacional.com/literatura",
				TopicLabel: "literature",
			},
			seedURL{
				URL:        "https://www.elnacional.com/musica",
				TopicLabel: "music",
			},
			seedURL{
				URL:        "https://www.elnacional.com/teatro",
				TopicLabel: "theater",
			},
			seedURL{
				URL:        "https://www.elnacional.com/opinion/",
				TopicLabel: "opinion",
			},
			seedURL{
				URL:        "https://www.elnacional.com/gastronomia",
				TopicLabel: "cooking",
			},
			seedURL{
				URL:        "https://www.elnacional.com/moda",
				TopicLabel: "fashion",
			},
			seedURL{
				URL:        "https://www.elnacional.com/viajes-turismo",
				TopicLabel: "travel",
			},
		},
	},
	"eldiario.com": allowableDomain{
		Country:      geo.CountryCodeVenezuela,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://eldiario.com/seccion/politica/",
				TopicLabel: "politics",
			},
			seedURL{
				URL:        "https://eldiario.com/seccion/economia/",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://eldiario.com/seccion/venezuela/",
				TopicLabel: "current-events-venezuela",
			},
			seedURL{
				URL:        "https://eldiario.com/seccion/mundo/",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://eldiario.com/seccion/deportes/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://eldiario.com/seccion/cultura/",
				TopicLabel: "culture",
			},
			seedURL{
				URL:        "https://eldiario.com/seccion/tecnologia/",
				TopicLabel: "technology",
			},
		},
	},
	"brecha.com.uy": allowableDomain{
		Country:      geo.CountryCodeUruguay,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://brecha.com.uy/category/politica/",
				TopicLabel: "politics",
			},
			seedURL{
				URL:        "https://brecha.com.uy/category/cultura/",
				TopicLabel: "culture",
			},
			seedURL{
				URL:        "https://brecha.com.uy/category/mundo/",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://brecha.com.uy/category/columnas-de-opinion/",
				TopicLabel: "opinion",
			},
		},
	},
	"ladiaria.com.uy": allowableDomain{
		Country:      geo.CountryCodeUruguay,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://ladiaria.com.uy/seccion/politica-nacional/",
				TopicLabel: "current-events-uruguay",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/seccion/politica-internacional/",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/opinion/",
				TopicLabel: "opinion",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/seccion/arte/",
				TopicLabel: "art",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/seccion/cine-tv-streaming/",
				TopicLabel: "film",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/seccion/teatro/",
				TopicLabel: "theater",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/seccion/musica/",
				TopicLabel: "music",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/seccion/videojuegos/",
				TopicLabel: "video-games",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/seccion/casa/",
				TopicLabel: "home",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/seccion/comer-y-beber/",
				TopicLabel: "cooking",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/seccion/tecnologia/",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/seccion/turismo/",
				TopicLabel: "travel",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/seccion/vida-saludable/",
				TopicLabel: "health",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/deporte/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/ciencia/",
				TopicLabel: "science",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/economia/",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://ladiaria.com.uy/salud/",
				TopicLabel: "health",
			},
		},
	},
	"ngenespanol.com": allowableDomain{
		Country:      geo.CountryCodeUnitedStates,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.ngenespanol.com/traveler/",
				TopicLabel: "travel",
			},
			seedURL{
				URL:        "https://www.ngenespanol.com/animales/",
				TopicLabel: "environment",
			},
			seedURL{
				URL:        "https://www.ngenespanol.com/el-mundo/",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.ngenespanol.com/el-espacio/",
				TopicLabel: "astronomy",
			},
			seedURL{
				URL:        "https://www.ngenespanol.com/ciencia/",
				TopicLabel: "science",
			},
		},
	},
	"primerahora.com": allowableDomain{
		Country:      geo.CountryCodePuertoRico,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		SeedURLs: []seedURL{
			seedURL{
				URL:        "https://www.primerahora.com/noticias/puerto-rico/",
				TopicLabel: "current-events-puerto-rico",
			},
			seedURL{
				URL:        "https://www.primerahora.com/noticias/policia-tribunales/",
				TopicLabel: "politics",
			},
			seedURL{
				URL:        "https://www.primerahora.com/noticias/gobierno-politica/",
				TopicLabel: "politics",
			},
			seedURL{
				URL:        "https://www.primerahora.com/noticias/consumo/",
				TopicLabel: "economy",
			},
			seedURL{
				URL:        "https://www.primerahora.com/noticias/estados-unidos/",
				TopicLabel: "current-events-united-states",
			},
			seedURL{
				URL:        "https://www.primerahora.com/noticias/mundo/",
				TopicLabel: "world-news",
			},
			seedURL{
				URL:        "https://www.primerahora.com/noticias/ciencia-tecnologia/",
				TopicLabel: "technology",
			},
			seedURL{
				URL:        "https://www.primerahora.com/entretenimiento/farandula/",
				TopicLabel: "entertainment",
			},
			seedURL{
				URL:        "https://www.primerahora.com/entretenimiento/musica/",
				TopicLabel: "music",
			},
			seedURL{
				URL:        "https://www.primerahora.com/entretenimiento/cine-tv/",
				TopicLabel: "film",
			},
			seedURL{
				URL:        "https://www.primerahora.com/entretenimiento/cultura-teatro/",
				TopicLabel: "theater",
			},
			seedURL{
				URL:        "https://www.primerahora.com/entretenimiento/otras/",
				TopicLabel: "entertainment",
			},
			seedURL{
				URL:        "https://www.primerahora.com/deportes/baloncesto/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.primerahora.com/deportes/beisbol/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.primerahora.com/deportes/boxeo/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.primerahora.com/deportes/voleibol/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.primerahora.com/deportes/hipismo/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.primerahora.com/deportes/otros/",
				TopicLabel: "sports",
			},
			seedURL{
				URL:        "https://www.primerahora.com/estilos-de-vida/ph-mas-saludable/",
				TopicLabel: "health",
			},
			seedURL{
				URL:        "https://www.primerahora.com/estilos-de-vida/cocina/",
				TopicLabel: "cooking",
			},
			seedURL{
				URL:        "https://www.primerahora.com/estilos-de-vida/moda-estilo/",
				TopicLabel: "fashion",
			},
		},
	},
}
