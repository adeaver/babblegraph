package ses

type Bounce struct {
	BounceID          string             `json:"feedbackId"`
	Type              BounceType         `json:"bounceType"`
	Subtype           string             `json:"bounceSubType"`
	BouncedRecipients []BouncedRecipient `json:"bouncedRecipients"`
	TimestampISO8601  string             `json:"timestamp"`
	RemoteMTAIP       string             `json:"remoteMtaIp"`
	ReportingMTA      string             `json:"reportingMta"`
}

type BouncedRecipient struct {
	EmailAddress   string  `json:"emailAddress"`
	Action         *string `json:"action,omitempty"`
	Status         *string `json:"status,omitempty"`
	DiagnosticCode *string `json:"diagnosticCode,omitempty"`
}

type BounceType string

const (
	BounceTypeUndetermined BounceType = "Undeteremined"
	BounceTypePermanent    BounceType = "Permanent"
	BounceTypeTransient    BounceType = "Transient"
)
