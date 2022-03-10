package ingesthtml

import (
	"babblegraph/model/content"
	"babblegraph/util/ptr"
	"babblegraph/util/testutils"
	"testing"
)

var (
	normalHTMLPage = `<html lang="es">
<head>
		<title>Page Title</title>
		<meta property="og:type" content="article" />
		<meta property="og:title" content="This is an Article" />
		<meta data-ue-u="og:title" content="Not an article" />
		<link rel="stylesheet" href="stylesheet.css" />
</head>
<body>
		<p>Some body text</p>
		<p>Text with a <a href="www.google.com">link</a></p>
		<p>Text with <strong>styling</strong> and a <span>span</span></p>
		<div>Text in a div</div>
		<script>
			some random javascript
		</script>
		<a href="/relative-link">relative link</a>
</body>`
	expectedParsedPage = ParsedHTMLPage{
		Links: []string{
			"www.google.com",
			"laprensa.com.ar/relative-link",
		},
		BodyText: "Some body text Text with a link Text with styling and a span Text in a div relative link",
		Language: ptr.String("es"),
		PageType: ptr.String("article"),
		Metadata: map[string]string{
			"og:type":  "article",
			"og:title": "This is an Article",
		},
	}
)

func TestParseHTML(t *testing.T) {
	parsed, err := parseHTML(parseHTMLInput{
		source: content.Source{
			URL: "laprensa.com.ar",
		},
		htmlStr: normalHTMLPage,
		cset:    "utf-8",
	})
	if err != nil {
		t.Errorf("Not expecting error, but got one: %s", err.Error())
		return
	}
	if err := testutils.CompareStringLists(expectedParsedPage.Links, parsed.Links); err != nil {
		t.Errorf("Error on links: %s", err.Error())
	}
	if err := testutils.CompareNullableString(parsed.Language, expectedParsedPage.Language); err != nil {
		t.Errorf("Error on Language: %s", err.Error())
	}
	if err := testutils.CompareNullableString(parsed.PageType, expectedParsedPage.PageType); err != nil {
		t.Errorf("Error on page type: %s", err.Error())
	}
	if err := testutils.CompareStringMap(parsed.Metadata, expectedParsedPage.Metadata); err != nil {
		t.Errorf("Error on metadata: %s", err.Error())
	}
}

func TestParseLDJSONPaywalledHTML(t *testing.T) {
	ldjsonPaywalledHTML := `<html lang="es">
<head>
		<title>Page Title</title>
		<meta property="og:type" content="article" />
		<meta property="og:title" content="This is an Article" />
		<meta data-ue-u="og:title" content="Not an article" />
		<link rel="stylesheet" href="stylesheet.css" />
        <script type="application/ld+json">
        {
            "@context": "http://schema.org",
            "@type":
                "OpinionNewsArticle",
                "mainEntityOfPage":{
                    "@type":"WebPage",
                    "@id":"https://www.elmundo.es/internacional/2021/05/20/60a53d22fc6c83a70e8b45c8.html"
                },
                "headline": "Rusia quiere mandar en el nuevo Ártico",
                "articleSection": "internacional",
                "datePublished": "2021-05-19T23:30:36Z",
                "dateModified": "2021-05-19T23:30:36Z",
                "image":{
                    "@type": "ImageObject",
                    "url": "https://phantom-elmundo.unidadeditorial.es/1fb5c7650df5663badbd70bcbdcd0b00/resize/1200/f/jpg/assets/multimedia/imagenes/2021/05/19/16214418327466.jpg",
                    "height": 800,
                    "width": 1200
                },
                "publisher": {
                    "@type": "NewsMediaOrganization",
                    "name": "El mundo",
                        "logo": {
                            "@type": "ImageObject",
                            "url": "https://e00-elmundo.uecdn.es/assets/desktop/master/img/iconos/elmundo.png",
                            "width": 204,
                            "height": 27
                        }
                },
                "description": "Rusia está reafirmando su posición como gran potencia mundial también en el Ártico, una zona que desde que Vladimir Putin llegó al Kremlin se ha convertido en un escenario...",
                "isAccessibleForFree": false,
                "hasPart":
                {
                    "@type": "WebPageElement",
                    "isAccessibleForFree": false,
                    "cssSelector": ".paywall"
                }
        }
        </script>
</head>
<body>
		<p>Some body text</p>
		<p>Text with a <a href="www.google.com">link</a></p>
		<p>Text with <strong>styling</strong> and a <span>span</span></p>
		<div>Text in a div</div>
		<script>
			some random javascript
		</script>
		<a href="/relative-link">relative link</a>
</body>`
	parsed, err := parseHTML(parseHTMLInput{
		source: content.Source{
			URL: "elmundo.es",
		},
		sourceFilter: &content.SourceFilter{
			UseLDJSONValidation: ptr.Bool(true),
		},
		htmlStr: ldjsonPaywalledHTML,
		cset:    "utf-8",
	})
	if err != nil {
		t.Errorf("Not expecting error, but got one: %s", err.Error())
		return
	}
	if !parsed.IsPaywalled {
		t.Errorf("Expected content to be paywalled, but isn't")
	}
}

func TestParseLDJSONNotPaywalledHTML(t *testing.T) {
	ldjsonPaywalledHTML := `<html lang="es">
<head>
		<title>Page Title</title>
		<meta property="og:type" content="article" />
		<meta property="og:title" content="This is an Article" />
		<meta data-ue-u="og:title" content="Not an article" />
		<link rel="stylesheet" href="stylesheet.css" />
        <script type="application/ld+json">
        {
            "@context": "http://schema.org",
            "@type":
                "OpinionNewsArticle",
                "mainEntityOfPage":{
                    "@type":"WebPage",
                    "@id":"https://www.elmundo.es/internacional/2021/05/20/60a53d22fc6c83a70e8b45c8.html"
                },
                "headline": "Rusia quiere mandar en el nuevo Ártico",
                "articleSection": "internacional",
                "datePublished": "2021-05-19T23:30:36Z",
                "dateModified": "2021-05-19T23:30:36Z",
                "image":{
                    "@type": "ImageObject",
                    "url": "https://phantom-elmundo.unidadeditorial.es/1fb5c7650df5663badbd70bcbdcd0b00/resize/1200/f/jpg/assets/multimedia/imagenes/2021/05/19/16214418327466.jpg",
                    "height": 800,
                    "width": 1200
                },
                "publisher": {
                    "@type": "NewsMediaOrganization",
                    "name": "El mundo",
                        "logo": {
                            "@type": "ImageObject",
                            "url": "https://e00-elmundo.uecdn.es/assets/desktop/master/img/iconos/elmundo.png",
                            "width": 204,
                            "height": 27
                        }
                },
                "description": "Rusia está reafirmando su posición como gran potencia mundial también en el Ártico, una zona que desde que Vladimir Putin llegó al Kremlin se ha convertido en un escenario...",
                "hasPart":
                {
                    "@type": "WebPageElement",
                    "isAccessibleForFree": false,
                    "cssSelector": ".paywall"
                }
        }
        </script>
</head>
<body>
		<p>Some body text</p>
		<p>Text with a <a href="www.google.com">link</a></p>
		<p>Text with <strong>styling</strong> and a <span>span</span></p>
		<div>Text in a div</div>
		<script>
			some random javascript
		</script>
		<a href="/relative-link">relative link</a>
</body>`
	parsed, err := parseHTML(parseHTMLInput{
		source: content.Source{
			URL: "elmundo.es",
		},
		sourceFilter: &content.SourceFilter{
			UseLDJSONValidation: ptr.Bool(true),
		},
		htmlStr: ldjsonPaywalledHTML,
		cset:    "utf-8",
	})
	if err != nil {
		t.Errorf("Not expecting error, but got one: %s", err.Error())
		return
	}
	if parsed.IsPaywalled {
		t.Errorf("Expected content not to be paywalled, but it is")
	}
}

func TestParseClassesNotPaywalledHTML(t *testing.T) {
	classesNotPaywalledHTML := `<html lang="es">
<head>
		<title>Page Title</title>
		<meta property="og:type" content="article" />
		<meta property="og:title" content="This is an Article" />
		<meta data-ue-u="og:title" content="Not an article" />
		<link rel="stylesheet" href="stylesheet.css" />
</head>
<body>
		<p>Some body text</p>
		<p>Text with a <a href="www.google.com">link</a></p>
		<p>Text with <strong>styling</strong> and a <span>span</span></p>
		<div>Text in a div</div>
		<script>
			some random javascript
		</script>
		<a href="/relative-link">relative link</a>
</body>`
	parsed, err := parseHTML(parseHTMLInput{
		source: content.Source{
			URL: "elespectador.com",
		},
		sourceFilter: &content.SourceFilter{
			PaywallClasses: []string{"premium_validation"},
		},
		htmlStr: classesNotPaywalledHTML,
		cset:    "utf-8",
	})
	if err != nil {
		t.Errorf("Not expecting error, but got one: %s", err.Error())
		return
	}
	if parsed.IsPaywalled {
		t.Errorf("Expected content not to be paywalled, but it is")
	}
}

func TestParseClassesPaywalledHTML(t *testing.T) {
	classesPaywalledHTML := `<html lang="es">
<head>
		<title>Page Title</title>
		<meta property="og:type" content="article" />
		<meta property="og:title" content="This is an Article" />
		<meta data-ue-u="og:title" content="Not an article" />
		<link rel="stylesheet" href="stylesheet.css" />
</head>
<body>
		<p>Some body text</p>
		<p>Text with a <a href="www.google.com">link</a></p>
		<p>Text with <strong>styling</strong> and a <span class="something something2 premium_validation something_else">span</span></p>
		<div>Text in a div</div>
		<script>
			some random javascript
		</script>
		<a href="/relative-link">relative link</a>
</body>`
	parsed, err := parseHTML(parseHTMLInput{
		source: content.Source{
			URL: "elespectador.com",
		},
		sourceFilter: &content.SourceFilter{
			PaywallClasses: []string{"premium_validation"},
		},
		htmlStr: classesPaywalledHTML,
		cset:    "utf-8",
	})
	if err != nil {
		t.Errorf("Not expecting error, but got one: %s", err.Error())
		return
	}
	if !parsed.IsPaywalled {
		t.Errorf("Expected content to be paywalled, but it is not")
	}
}

func TestParseLDJSONPaywalledHTMLWithString(t *testing.T) {
	ldjsonPaywalledHTML := `<html lang="es">
<head>
		<title>Page Title</title>
		<meta property="og:type" content="article" />
		<meta property="og:title" content="This is an Article" />
		<meta data-ue-u="og:title" content="Not an article" />
		<link rel="stylesheet" href="stylesheet.css" />
        <script type="application/ld+json">
        {
            "@context": "http://schema.org",
            "@type":
                "OpinionNewsArticle",
                "mainEntityOfPage":{
                    "@type":"WebPage",
                    "@id":"https://www.elmundo.es/internacional/2021/05/20/60a53d22fc6c83a70e8b45c8.html"
                },
                "headline": "Rusia quiere mandar en el nuevo Ártico",
                "articleSection": "internacional",
                "datePublished": "2021-05-19T23:30:36Z",
                "dateModified": "2021-05-19T23:30:36Z",
                "image":{
                    "@type": "ImageObject",
                    "url": "https://phantom-elmundo.unidadeditorial.es/1fb5c7650df5663badbd70bcbdcd0b00/resize/1200/f/jpg/assets/multimedia/imagenes/2021/05/19/16214418327466.jpg",
                    "height": 800,
                    "width": 1200
                },
                "publisher": {
                    "@type": "NewsMediaOrganization",
                    "name": "El mundo",
                        "logo": {
                            "@type": "ImageObject",
                            "url": "https://e00-elmundo.uecdn.es/assets/desktop/master/img/iconos/elmundo.png",
                            "width": 204,
                            "height": 27
                        }
                },
                "description": "Rusia está reafirmando su posición como gran potencia mundial también en el Ártico, una zona que desde que Vladimir Putin llegó al Kremlin se ha convertido en un escenario...",
                "isAccessibleForFree": false,
                "hasPart":
                {
                    "@type": "WebPageElement",
                    "isAccessibleForFree": "False",
                    "cssSelector": ".paywall"
                }
        }
        </script>
</head>
<body>
		<p>Some body text</p>
		<p>Text with a <a href="www.google.com">link</a></p>
		<p>Text with <strong>styling</strong> and a <span>span</span></p>
		<div>Text in a div</div>
		<script>
			some random javascript
		</script>
		<a href="/relative-link">relative link</a>
</body>`
	parsed, err := parseHTML(parseHTMLInput{
		source: content.Source{
			URL: "elmundo.es",
		},
		sourceFilter: &content.SourceFilter{
			UseLDJSONValidation: ptr.Bool(true),
		},
		htmlStr: ldjsonPaywalledHTML,
		cset:    "utf-8",
	})
	if err != nil {
		t.Errorf("Not expecting error, but got one: %s", err.Error())
		return
	}
	if !parsed.IsPaywalled {
		t.Errorf("Expected content to be paywalled, but isn't")
	}
}

func TestParseIDsPaywalledHTML(t *testing.T) {
	idsPaywalledHTML := `<html lang="es">
<head>
		<title>Page Title</title>
		<meta property="og:type" content="article" />
		<meta property="og:title" content="This is an Article" />
		<meta data-ue-u="og:title" content="Not an article" />
		<link rel="stylesheet" href="stylesheet.css" />
</head>
<body>
		<p>Some body text</p>
		<p>Text with a <a href="www.google.com">link</a></p>
		<p id="is_c9_article">Text with <strong>styling</strong> and a <span>span</span></p>
		<div>Text in a div</div>
		<script>
			some random javascript
		</script>
		<a href="/relative-link">relative link</a>
</body>`
	parsed, err := parseHTML(parseHTMLInput{
		source: content.Source{
			URL: "yucatan.com.mx",
		},
		sourceFilter: &content.SourceFilter{
			PaywallIDs: []string{"is_c9_article"},
		},
		htmlStr: idsPaywalledHTML,
		cset:    "utf-8",
	})
	if err != nil {
		t.Errorf("Not expecting error, but got one: %s", err.Error())
		return
	}
	if !parsed.IsPaywalled {
		t.Errorf("Expected content to be paywalled, but it is not")
	}
}

func TestParseNotIDsPaywalledHTML(t *testing.T) {
	idsNotPaywalledHTML := `<html lang="es">
<head>
		<title>Page Title</title>
		<meta property="og:type" content="article" />
		<meta property="og:title" content="This is an Article" />
		<meta data-ue-u="og:title" content="Not an article" />
		<link rel="stylesheet" href="stylesheet.css" />
</head>
<body>
		<p>Some body text</p>
		<p id="not-a-paywall">Text with a <a href="www.google.com">link</a></p>
		<p>Text with <strong>styling</strong> and a <span>span</span></p>
		<div>Text in a div</div>
		<script>
			some random javascript
		</script>
		<a href="/relative-link">relative link</a>
</body>`
	parsed, err := parseHTML(parseHTMLInput{
		source: content.Source{
			URL: "yucatan.com.mx",
		},
		sourceFilter: &content.SourceFilter{
			PaywallIDs: []string{"is_c9_article"},
		},
		htmlStr: idsNotPaywalledHTML,
		cset:    "utf-8",
	})
	if err != nil {
		t.Errorf("Not expecting error, but got one: %s", err.Error())
		return
	}
	if parsed.IsPaywalled {
		t.Errorf("Expected content not to be paywalled, but it is")
	}
}
