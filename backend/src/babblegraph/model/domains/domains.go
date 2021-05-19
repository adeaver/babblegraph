package domains

import (
	"babblegraph/util/geo"
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
)

var allowableDomains = []AllowableDomain{
	{
		Domain:       "ambito.com",
		Country:      geo.CountryCodeArgentina,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "cronista.com",
		Country:      geo.CountryCodeArgentina,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "laprensa.com.ar",
		Country:      geo.CountryCodeArgentina,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "emol.com",
		Country:      geo.CountryCodeChile,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:                      "latercera.com",
		Country:                     geo.CountryCodeChile,
		LanguageCode:                wordsmith.LanguageCodeSpanish,
		NumberOfMonthlyFreeArticles: ptr.Int64(15),
	}, {
		Domain:       "lacuarta.cl",
		Country:      geo.CountryCodeChile,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "elbogotano.com.co",
		Country:      geo.CountryCodeColombia,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "elespectador.com",
		Country:      geo.CountryCodeColombia,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		PaywallValidation: &PaywallValidation{
			PaywallClasses: []string{"premium_validation"},
		},
	}, {
		Domain:                      "elsalvador.com",
		Country:                     geo.CountryCodeElSalvador,
		LanguageCode:                wordsmith.LanguageCodeSpanish,
		NumberOfMonthlyFreeArticles: ptr.Int64(20),
	}, {
		Domain:       "elmundo.sv",
		Country:      geo.CountryCodeElSalvador,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "dca.gob.gt",
		Country:      geo.CountryCodeGuatemala,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "prensalibre.com",
		Country:      geo.CountryCodeGuatemala,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		PaywallValidation: &PaywallValidation{
			PaywallClasses: []string{"pl_plus", "type-pl_plus"},
		},
	}, {
		Domain:       "laprensa.hn",
		Country:      geo.CountryCodeHonduras,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "heraldodemexico.com.mx",
		Country:      geo.CountryCodeMexico,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "eleconomista.com.mx",
		Country:      geo.CountryCodeMexico,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "yucatan.com.mx",
		Country:      geo.CountryCodeMexico,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		PaywallValidation: &PaywallValidation{
			PaywallClasses: []string{"tag-central-9"},
		},
	}, {
		Domain:       "informador.mx",
		Country:      geo.CountryCodeMexico,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:                      "laprensa.com.ni",
		Country:                     geo.CountryCodeNicaragua,
		LanguageCode:                wordsmith.LanguageCodeSpanish,
		NumberOfMonthlyFreeArticles: ptr.Int64(5),
		PaywallValidation: &PaywallValidation{
			PaywallClasses: []string{"tag-exclusivo"},
		},
	}, {
		Domain:       "lajornadanet.com",
		Country:      geo.CountryCodeNicaragua,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "elsiglo.com",
		Country:      geo.CountryCodePanama,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "larepublica.pe",
		Country:      geo.CountryCodePeru,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "abc.com.py",
		Country:      geo.CountryCodeParaguay,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "ultimahora.com",
		Country:      geo.CountryCodeParaguay,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "abc.es",
		Country:      geo.CountryCodeSpain,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		PaywallValidation: &PaywallValidation{
			UseLDJSONValidation: &struct{}{},
		},
	}, {
		Domain:       "elmundo.es",
		Country:      geo.CountryCodeSpain,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		PaywallValidation: &PaywallValidation{
			UseLDJSONValidation: &struct{}{},
		},
	}, {
		Domain:                      "elpais.com",
		Country:                     geo.CountryCodeSpain,
		LanguageCode:                wordsmith.LanguageCodeSpanish,
		NumberOfMonthlyFreeArticles: ptr.Int64(10),
	}, {
		Domain:                      "elcomercio.es",
		Country:                     geo.CountryCodeSpain,
		LanguageCode:                wordsmith.LanguageCodeSpanish,
		NumberOfMonthlyFreeArticles: ptr.Int64(10),
		PaywallValidation: &PaywallValidation{
			UseLDJSONValidation: &struct{}{},
		},
	}, {
		Domain:       "diariolasamericas.com",
		Country:      geo.CountryCodeUnitedStates,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:                      "laopinion.com",
		Country:                     geo.CountryCodeUnitedStates,
		LanguageCode:                wordsmith.LanguageCodeSpanish,
		NumberOfMonthlyFreeArticles: ptr.Int64(20),
	}, {
		Domain:                      "eldiariony.com",
		Country:                     geo.CountryCodeUnitedStates,
		LanguageCode:                wordsmith.LanguageCodeSpanish,
		NumberOfMonthlyFreeArticles: ptr.Int64(20),
	}, {
		Domain:       "hoylosangeles.com",
		Country:      geo.CountryCodeUnitedStates,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "elnacional.com",
		Country:      geo.CountryCodeVenezuela,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "eldiario.com",
		Country:      geo.CountryCodeVenezuela,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "brecha.com.uy",
		Country:      geo.CountryCodeUruguay,
		LanguageCode: wordsmith.LanguageCodeSpanish,
		PaywallValidation: &PaywallValidation{
			UseLDJSONValidation: &struct{}{},
		},
	}, {
		Domain:                      "elpais.com.uy",
		Country:                     geo.CountryCodeUruguay,
		LanguageCode:                wordsmith.LanguageCodeSpanish,
		NumberOfMonthlyFreeArticles: ptr.Int64(10),
	}, {
		Domain:                      "ladiaria.com.uy",
		Country:                     geo.CountryCodeUruguay,
		LanguageCode:                wordsmith.LanguageCodeSpanish,
		NumberOfMonthlyFreeArticles: ptr.Int64(10),
	}, {
		Domain:       "ngenespanol.com",
		Country:      geo.CountryCodeUnitedStates,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	}, {
		Domain:       "primerahora.com",
		Country:      geo.CountryCodePuertoRico,
		LanguageCode: wordsmith.LanguageCodeSpanish,
	},
}
