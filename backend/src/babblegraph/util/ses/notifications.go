package ses

// defined in https://docs.aws.amazon.com/ses/latest/DeveloperGuide/notification-contents.html
type Notification struct {
	MessageID        SESMessageID `json:"MessageId,omitempty"`
	Token            *string      `json:"Token,omitempty"`
	TopicARN         string       `json:"TopicArn,omitempty"`
	Message          *string      `json:"Message,omitempty"`
	TimestampISO8601 string       `json:"Timestamp,omitempty"`
	SignatureVersion string       `json:"SignatureVersion,omitempty"`
	Signature        string       `json:"Signature,omitempty"`
	SigningCertURL   string       `json:"SigningCertURL,omitempty"`

	SubscribeURL   *string `json:"SubscribeURL,omitempty"`
	UnsubscribeURL *string `json:"UnsubscribeURL,omitempty"`
}

type SESMessageID string

type NotificationType string

const (
	NotificationTypeBounce                   NotificationType = "Bounce"
	NotificationTypeComplaint                NotificationType = "Complaint"
	NotificationTypeDelivery                 NotificationType = "Delivery"
	NotificationTypeSubscriptionConfirmation NotificationType = "SubscriptionConfirmation"
)

type NotificationBody struct {
	Type      NotificationType `json:"notificationType"`
	Mail      Mail             `json:"mail"`
	Bounce    *Bounce          `json:"bounce,omitempty"`
	Complaint *Complaint       `json:"complaint,omitempty"`
	Delivery  *Delivery        `json:"delivery,omitempty"`
}

type Mail struct {
	TimestampISO8601    string         `json:"timestamp"`
	OriginalMessageID   string         `json:"messageId"`
	FromAddress         string         `json:"source"`
	FromAddressARN      string         `json:"sourceArn"`
	FromIP              string         `json:"sourceIp"`
	SendingAccountID    string         `json:"sendingAccountId"`
	Destination         []string       `json:"destination"`
	AreHeadersTruncated *bool          `json:"headersTruncated,omitempty"`
	Headers             []Header       `json:"headers,omitempty"`
	CommonHeaders       *CommonHeaders `json:"commonHeaders,omitempty"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CommonHeaders struct {
	From      []string `json:"from"`
	Date      string   `json:"date"`
	To        []string `json:"to"`
	MessageID string   `json:"messageId"`
	Subject   string   `json:"subject"`
}
