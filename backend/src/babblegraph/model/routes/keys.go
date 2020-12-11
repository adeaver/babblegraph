package routes

type RouteEncryptionKey string

const (
	SubscriptionManagementRouteEncryptionKey RouteEncryptionKey = "subscription-management"
)

func (r RouteEncryptionKey) Str() string {
	return string(r)
}
