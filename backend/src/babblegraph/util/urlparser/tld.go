package urlparser

// The URL parser is only accepting combinations of these TLDs
// it is far from an exhaustive list, but we want to err on the side
// of too strict than too loose since the penalty for missing content
// is significantly lower than the penalty for delivering bad or duplicate content
var validTLDs = map[string]bool{
	".com": true,
	".net": true,
	".org": true,
	".gov": true,
	".edu": true,
	".nom": true,
	".gob": true,
	".ngo": true,
	".mil": true,
	".inf": true,

	".es": true, // Spain
	".mx": true, // Mexico
	".co": true, // Colombia
	".ar": true, // Argentina
	".pe": true, // Peru
	".ve": true, // Venezuela
	".cl": true, // Chile
	".ec": true, // Ecuador
	".gt": true, // Guatemala
	".cu": true, // Cuba
	".bo": true, // Bolivia
	".do": true, // Dominican Republic
	".hn": true, // Honduras
	".py": true, // Paraguay
	".sv": true, // El Salvador
	".ni": true, // Nicaragua
	".cr": true, // Costa Rica
	".pa": true, // Panama
	".uy": true, // Uruguay
	".gq": true, // Equatorial Guinea
	".pr": true, // Puerto Rico
	".bz": true, // Belize

	".musica": true, // Random second level domain for Argentina
}
