package clienterror

type Error string

const (
	// General Errors
	ErrorInvalidToken    Error = "invalid-token"
	ErrorNoAuth          Error = "no-auth"
	ErrorIncorrectKey    Error = "incorrect-key"
	ErrorLowCaptchaScore Error = "low-score"
	ErrorRequiresUpgrade Error = "requires-upgrade"

	// Common Error Types
	ErrorInvalidLanguageCode Error = "invalid-language"
	ErrorInvalidEmailAddress Error = "invalid-email-address"
)

func (e Error) Ptr() *Error {
	return &e
}
