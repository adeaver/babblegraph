package blog

type Content []ContentNode

type ContentNode struct {
	Type ContentNodeType `json:"type"`
	Body interface{}     `json:"body"`
}

func (c Content) RemoveNode(idx int) {
	c = append(c[:idx], c[idx+1:]...)
}

type ContentNodeType string

const (
	ContentNodeTypeHeading   ContentNodeType = "heading"
	ContentNodeTypeParagraph ContentNodeType = "paragraph"
)

type Heading struct {
	Text string `json:"text"`
}

func (c Content) AddHeading(h Heading) {
	c = append(c, ContentNode{
		Type: ContentNodeTypeHeading,
		Body: h,
	})
}

type Paragraph struct {
	Text string `json:"text"`
}

func (c Content) AddParagraph(p Paragraph) {
	c = append(c, ContentNode{
		Type: ContentNodeTypeParagraph,
		Body: p,
	})
}
