package domain

type OrderStatus string

const (
	OrderPlaced     OrderStatus = "placed"
	OrderDispatched OrderStatus = "dispatched"
	OrderCompleted  OrderStatus = "completed"
	OrderCancelled  OrderStatus = "cancelled"
)
