package routes

type RouteEncryptionKey string

const (
	SubscriptionManagementRouteEncryptionKey RouteEncryptionKey = "subscription-management"
	UnsubscribeRouteEncryptionKey            RouteEncryptionKey = "unsubscribe"
	EmailOpenedKey                           RouteEncryptionKey = "email-opened"
	UserVerificationKey                      RouteEncryptionKey = "user-verification"
	WordReinforcementKey                     RouteEncryptionKey = "word-reinforcement"
)

func (r RouteEncryptionKey) Str() string {
	return string(r)
}
