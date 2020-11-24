package opengraph

type Tag string

const (
	TitleTag           Tag = "og:title"
	ImageURLTag        Tag = "og:image"
	TypeTag            Tag = "og:type"
	URLTag             Tag = "og:url"
	DescriptionTag     Tag = "og:description"
	PublicationTimeTag Tag = "og:article:published_time"
)

func (t Tag) Ptr() *Tag {
	return &t
}

func (t Tag) Str() string {
	return string(t)
}
