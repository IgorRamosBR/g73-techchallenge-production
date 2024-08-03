package gateways

type OrderRepository interface {
	SaveOrder() error
}
