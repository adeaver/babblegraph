package ses

type Delivery struct {
	TimestampISO8601           string   `json:"timestamp"`
	ProcessingTimeMilliseconds int64    `json:"processingTimeMillis"`
	Recipients                 []string `json:"recipients"`
	SMTPResponse               string   `json:"smtpResponse"`
	ReportingMTA               string   `json:"reportingMTA"`
	ReportingMTAIP             string   `json:"remoteMtaIp"`
}
