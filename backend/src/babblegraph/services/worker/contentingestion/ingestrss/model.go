package ingestrss

import "encoding/xml"

type podcastRSSFeed struct {
	XMLName xml.Name          `xml:"rss"`
	Channel podcastRSSChannel `xml:"channel"`
}

type podcastRSSChannel struct {
	XMLName        xml.Name              `xml:"channel"`
	AtomLink       string                `xml:"atom:link"`
	Title          string                `xml:"title"`
	Link           string                `xml:"link"`
	Language       string                `xml:"language"`
	Copyright      string                `xml:"copyright"`
	Description    string                `xml:"description"`
	Image          podcastRSSImage       `xml:"image"`
	IsExplicit     podcastExplicitity    `xml:"itunes:explicit"`
	PodcastType    podcastType           `xml:"itunes:type"`
	Subtitle       string                `xml:"itunes:subtitle"`
	Author         string                `xml:"itunes:author"`
	Summary        string                `xml:"itunes:summary"`
	ContentEncoded podcastEncodedContent `xml:"content:encoded"`
	Owner          podcastOwner          `xml:"itunes:owner"`
	ITunesImage    podcastITunesImage    `xml:"itunes:image"`
	Categories     []podcastCategory     `xml:"itunes:category"`
	RSSFeedURL     string                `xml:"itunes:new-feed-url"`
	Episodes       []podcastEpisode      `xml:"item"`
}

type podcastRSSImage struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:"url"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
}

type podcastExplicitity string

const (
	podcastExplicitityYes   podcastExplicitity = "yes"
	podcastExplicitityTrue  podcastExplicitity = "true"
	podcastExplicitityNo    podcastExplicitity = "no"
	podcastExplicitityFalse podcastExplicitity = "false"
)

type podcastType string

const (
	podcastTypeEpisodic podcastType = "episodic"
)

type podcastEncodedContent struct {
	XMLName xml.Name `xml:"content:encoded"`
	Value   string   `xml:",cdata"`
}

type podcastOwner struct {
	XMLName xml.Name `xml:"itunes:owner"`
	Name    string   `xml:"itunes:name"`
	Email   string   `xml:"itunes:email"`
}

type podcastITunesImage struct {
	XMLName xml.Name `xml:"itunes:image"`
	URL     string   `xml:"href,attr"`
}

type podcastCategory struct {
	XMLName xml.Name `xml:"itunes:category"`
	Name    string   `xml:"text,attr"`
}

type podcastEpisode struct {
	XMLName         xml.Name              `xml:"item"`
	Title           string                `xml:"title"`
	Description     string                `xml:"description"`
	PublicationDate string                `xml:"pubDate"`
	EpisodeType     podcastEpisodeType    `xml:"itunes:episodeType"`
	Author          string                `xml:"author"`
	Subtitle        string                `xml:"subtitle"`
	Summary         string                `xml:"itunes:summary"`
	ContentEncoded  podcastEncodedContent `xml:"content:encoded"`
	Duration        string                `xml:"duration"`
	IsExplicit      podcastExplicitity    `xml:"itunes:explicit"`
	ID              string                `xml:"guid"`
	AudioData       podcastEnclosure      `xml:"enclosure"`
}

type podcastEpisodeType string

const (
	podcastEpisodeTypeFull podcastEpisodeType = "full"
)

type podcastEnclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}
