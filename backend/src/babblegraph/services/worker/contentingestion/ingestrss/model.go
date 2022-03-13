package ingestrss

import (
	"babblegraph/util/ctx"
	"encoding/xml"
	"strings"
)

type podcastRSSFeed struct {
	XMLName xml.Name          `xml:"rss"`
	Channel PodcastRSSChannel `xml:"channel"`
}

type PodcastRSSChannel struct {
	XMLName        xml.Name              `xml:"channel"`
	AtomLink       string                `xml:"atom:link"`
	Title          string                `xml:"title"`
	Link           string                `xml:"link"`
	Language       string                `xml:"language"`
	Copyright      string                `xml:"copyright"`
	Description    string                `xml:"description"`
	Image          PodcastRSSImage       `xml:"image"`
	IsExplicit     PodcastExplicitity    `xml:"itunes:explicit"`
	PodcastType    PodcastType           `xml:"itunes:type"`
	Subtitle       string                `xml:"itunes:subtitle"`
	Author         string                `xml:"itunes:author"`
	Summary        string                `xml:"itunes:summary"`
	ContentEncoded PodcastEncodedContent `xml:"content:encoded"`
	Owner          PodcastOwner          `xml:"itunes:owner"`
	ITunesImage    PodcastITunesImage    `xml:"itunes:image"`
	Categories     []PodcastCategory     `xml:"itunes:category"`
	RSSFeedURL     string                `xml:"itunes:new-feed-url"`
	Episodes       []PodcastEpisode      `xml:"item"`
}

type PodcastRSSImage struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:"url"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
}

type PodcastExplicitity string

const (
	PodcastExplicitityYes   PodcastExplicitity = "yes"
	PodcastExplicitityTrue  PodcastExplicitity = "true"
	PodcastExplicitityNo    PodcastExplicitity = "no"
	PodcastExplicitityFalse PodcastExplicitity = "false"
)

func (p PodcastExplicitity) Str() string {
	return string(p)
}

func (p PodcastExplicitity) ToBool(c ctx.LogContext) bool {
	switch strings.ToLower(p.Str()) {
	case PodcastExplicitityYes.Str(),
		PodcastExplicitityTrue.Str():
		return true
	case PodcastExplicitityNo.Str(),
		PodcastExplicitityFalse.Str():
		return false
	default:
		c.Warnf("Found unrecognized podcast explicit label %s. Assuming explicit", p.Str())
		return true
	}
}

type PodcastType string

const (
	PodcastTypeEpisodic PodcastType = "episodic"
)

type PodcastEncodedContent struct {
	XMLName xml.Name `xml:"content:encoded"`
	Value   string   `xml:",cdata"`
}

type PodcastOwner struct {
	XMLName xml.Name `xml:"itunes:owner"`
	Name    string   `xml:"itunes:name"`
	Email   string   `xml:"itunes:email"`
}

type PodcastITunesImage struct {
	XMLName xml.Name `xml:"itunes:image"`
	URL     string   `xml:"href,attr"`
}

type PodcastCategory struct {
	XMLName xml.Name `xml:"itunes:category"`
	Name    string   `xml:"text,attr"`
}

type PodcastEpisode struct {
	XMLName         xml.Name              `xml:"item"`
	Title           string                `xml:"title"`
	Description     string                `xml:"description"`
	PublicationDate string                `xml:"pubDate"`
	EpisodeType     PodcastEpisodeType    `xml:"itunes:episodeType"`
	Author          string                `xml:"author"`
	Subtitle        string                `xml:"subtitle"`
	Summary         string                `xml:"itunes:summary"`
	ContentEncoded  PodcastEncodedContent `xml:"content:encoded"`
	Duration        string                `xml:"duration"`
	IsExplicit      PodcastExplicitity    `xml:"itunes:explicit"`
	ID              string                `xml:"guid"`
	AudioData       PodcastEnclosure      `xml:"enclosure"`
}

type PodcastEpisodeType string

const (
	PodcastEpisodeTypeFull PodcastEpisodeType = "full"
)

func (p PodcastEpisodeType) Str() string {
	return string(p)
}

type PodcastEnclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}
