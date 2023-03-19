package serializer

type AddProductToOrderRequest struct {
	ProductID string `json:"product_id"`
}

type UpdateOrderRequest struct {
	DispatchDate string `json:"dispatch_date,omitempty"`
	OrderStatus  string `json:"order_status,omitempty"`
}
