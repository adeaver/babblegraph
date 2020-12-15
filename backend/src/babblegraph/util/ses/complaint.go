package ses

type Complaint struct {
	ComplaintID             string                `json:"feedbackId"`
	Type                    *ComplaintType        `json:"complaintFeedbackType,omitempty"`
	Subtype                 *string               `json:"complaintSubType,omitempty"`
	ComplainedRecipients    []ComplainedRecipient `json:"complainedRecipients"`
	TimestampISO8601        string                `json:"timestamp"`
	UserAgent               *string               `json:"userAgent,omitempty"`
	ArrivalTimestampISO8601 string                `json:"arrivalDate"`
}

type ComplainedRecipient struct {
	EmailAddress string `json:"emailAddress"`
}

type ComplaintType string

const (
	ComplaintTypeAbuse       ComplaintType = "abuse"
	ComplaintTypeAuthFailure ComplaintType = "auth-failure"
	ComplaintTypeFraud       ComplaintType = "fraud"
	ComplaintTypeNotSpam     ComplaintType = "not-spam"
	ComplaintTypeOther       ComplaintType = "other"
	ComplaintTypeVirus       ComplaintType = "virus"
)
