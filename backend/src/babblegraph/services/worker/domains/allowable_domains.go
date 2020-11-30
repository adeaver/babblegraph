package domains

var allowableDomains = map[string]bool{
	// Financial news, Argentina
	"ambito.com": true,

	// Daily news, Argentina
	"cronista.com":    true,
	"laprensa.com.ar": true,

	// Daily news, Chile
	"emol.com":      true,
	"latercera.com": true,

	// Celebrity news, Chile
	"lacuarta.cl": true,

	// Daily news, Colombia
	"elcolombiano.com":  true,
	"elbogotano.com.co": true,
	"elespectador.com":  true,

	// Daily news, Costa Rica
	"nacion.com": true,

	// Daily news, El Salvador
	"elsalvador.com": true,
	"elmundo.sv":     true,

	// Daily news, Guatemala
	"dca.gob.gt":      true,
	"prensalibre.com": true,

	// Daily news, Honduras
	"laprensa.hn": true,

	// Daily news, Mexico
	"jornada.com.mx":         true,
	"heraldodemexico.com.mx": true,

	// Financial news, Mexico
	"eleconomista.com.mx": true,

	// Daily news, Nicaragua
	"laprensa.com.ni":  true,
	"lajornadanet.com": true,

	// Daily news, Panama
	"elsiglo.com": true,

	// Daily news, Peru
	"larepublica.pe": true,

	// Daily news, Paraguay
	"abc.com.py":     true,
	"ultimahora.com": true,

	// Daily news, Spain
	"abc.es":        true,
	"elmundo.es":    true,
	"elpais.com":    true,
	"elcomercio.es": true,

	// Daily news, United States
	"diariolasamericas.com": true,
	"laopinion.com":         true,
	"eldiariony.com":        true,
	"hoylosangeles.com":     true,

	// Daily news, Venezuela
	"elnacional.com": true,
	"eldiario.com":   true,

	// Daily news, Uruguay
	"brecha.com.uy":   true,
	"elpais.com.uy":   true,
	"ladiaria.com.uy": true,

	// Daily news, Puerto Rico
	"periodismoinvestigativo.com": true,
}
