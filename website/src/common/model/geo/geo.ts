export enum CountryCode {
	Argentina = "AR",
	Chile = "CL",
	Colombia = "CO",
	CostaRica = "CR",
	ElSalvador = "SV",
	Guatemala = "GT",
	Honduras = "HN",
	Mexico = "MX",
	Nicaragua = "NI",
	Panama = "PA",
	Paraguay = "PY",
	Peru = "PE",
	PuertoRico = "PR",
	Spain = "ES",
	UnitedStates = "US",
	Uruguay = "UY",
	Venezuela = "VE",
}

export const getEnglishNameForCountryCode = (code: CountryCode) => {
    switch (code) {
	case CountryCode.Argentina:
        return "Argentina";
	case CountryCode.Chile:
        return "Chile";
	case CountryCode.Colombia:
        return "Colombia";
	case CountryCode.CostaRica:
        return "Costa Rica";
	case CountryCode.ElSalvador:
        return "El Salvador";
	case CountryCode.Guatemala:
        return "Guatemala";
	case CountryCode.Honduras:
        return "Honduras";
	case CountryCode.Mexico:
        return "Mexico";
	case CountryCode.Nicaragua:
        return "Nicaragua";
	case CountryCode.Panama:
        return "Panama";
	case CountryCode.Paraguay:
        return "Paraguay";
	case CountryCode.Peru:
        return "Peru";
	case CountryCode.PuertoRico:
        return "Puerto Rico";
	case CountryCode.Spain:
        return "Spain";
	case CountryCode.UnitedStates:
        return "United States";
	case CountryCode.Uruguay:
        return "Uruguay";
	case CountryCode.Venezuela:
        return "Venezuela";
    default:
        throw new Error(`Unrecognized country code: ${code}`);
    }
}
