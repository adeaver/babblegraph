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
	},
}
