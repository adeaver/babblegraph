# Allowable Domains

In order to ensure better quality content, I'm limiting ingest to verified domains.

The current verification process for a domain:
- Publication is reputable - I tried to avoid publications controlled by authoritarian governments or special interest groups
- No populist or anarchist news sources, however, conservative news publications are okay since conservative means slightly different things in different countries and in theory conservatism doesn't have to be a bad thing.
- No paywalls, although as of now, two domains on this list have monthly article caps (ladiaria.com.uy and elpais.com); if this becomes an issue, they will be removed
- Main webpage either has no og:type in its metadata, or a type of website
- The sections of the website have og:type of website or no og:type in their  metadata
- Articles have og:type of article.

## Plan to remove restrictions on metadata

In order to not have a restriction on metadata, I'll likely need a logistic regression model. The following domains were considered as part of the original group, but did not meet metadata requirements

- https://www.eldiario.net/ (Bolivia)
- https://www.diariolibre.com/ (Dominican Republic)
- http://www.elcomercio.com/ (Ecuador)
- http://www.cronica.com.mx/noticias.php (Mexico)
- http://laestrella.com.pa/ (Panama)
- http://www.prensa.com/ (Panama)
- http://www.eluniversal.com/ (Venezuela)

The requirement that sections be labeled with og:type website is not super strict as of 11/16/2020.

## Removed domains (with reason)

- https://jornada.com.mx/ (bad formatting - almost 15,000 articles had multiple content topics)
- https://www.elcolombiano.com/ (aggressive paywall)
