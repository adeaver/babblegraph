package geo

import (
	"fmt"
	"strings"
)

type CountryCode string

const (
	CountryCodeArgentina    CountryCode = "AR"
	CountryCodeChile        CountryCode = "CL"
	CountryCodeColombia     CountryCode = "CO"
	CountryCodeCostaRica    CountryCode = "CR"
	CountryCodeElSalvador   CountryCode = "SV"
	CountryCodeGuatemala    CountryCode = "GT"
	CountryCodeHonduras     CountryCode = "HN"
	CountryCodeMexico       CountryCode = "MX"
	CountryCodeNicaragua    CountryCode = "NI"
	CountryCodePanama       CountryCode = "PA"
	CountryCodeParaguay     CountryCode = "PY"
	CountryCodePeru         CountryCode = "PE"
	CountryCodePuertoRico   CountryCode = "PR"
	CountryCodeSpain        CountryCode = "ES"
	CountryCodeUnitedStates CountryCode = "US"
	CountryCodeUruguay      CountryCode = "UY"
	CountryCodeVenezuela    CountryCode = "VE"
)

func (c CountryCode) Str() string {
	return string(c)
}

func (c CountryCode) Ptr() *CountryCode {
	return &c
}

func GetCountryCodeFromString(c string) (*CountryCode, error) {
	switch strings.ToUpper(c) {
	case CountryCodeArgentina.Str():
		return CountryCodeArgentina.Ptr(), nil
	case CountryCodeChile.Str():
		return CountryCodeChile.Ptr(), nil
	case CountryCodeColombia.Str():
		return CountryCodeColombia.Ptr(), nil
	case CountryCodeCostaRica.Str():
		return CountryCodeCostaRica.Ptr(), nil
	case CountryCodeElSalvador.Str():
		return CountryCodeElSalvador.Ptr(), nil
	case CountryCodeGuatemala.Str():
		return CountryCodeGuatemala.Ptr(), nil
	case CountryCodeHonduras.Str():
		return CountryCodeHonduras.Ptr(), nil
	case CountryCodeMexico.Str():
		return CountryCodeMexico.Ptr(), nil
	case CountryCodeNicaragua.Str():
		return CountryCodeNicaragua.Ptr(), nil
	case CountryCodePanama.Str():
		return CountryCodePanama.Ptr(), nil
	case CountryCodeParaguay.Str():
		return CountryCodeParaguay.Ptr(), nil
	case CountryCodePeru.Str():
		return CountryCodePeru.Ptr(), nil
	case CountryCodePuertoRico.Str():
		return CountryCodePuertoRico.Ptr(), nil
	case CountryCodeSpain.Str():
		return CountryCodeSpain.Ptr(), nil
	case CountryCodeUnitedStates.Str():
		return CountryCodeUnitedStates.Ptr(), nil
	case CountryCodeUruguay.Str():
		return CountryCodeUruguay.Ptr(), nil
	case CountryCodeVenezuela.Str():
		return CountryCodeVenezuela.Ptr(), nil
	default:
		return nil, fmt.Errorf("Unsupported country code: %s", c)
	}
}
