package routes

type RouteEncryptionKey string

const (
	SubscriptionManagementRouteEncryptionKey RouteEncryptionKey = "subscription-management"
	UnsubscribeRouteEncryptionKey            RouteEncryptionKey = "unsubscribe"
)

func (r RouteEncryptionKey) Str() string {
	return string(r)
}
