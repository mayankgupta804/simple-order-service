package domain

import (
	"fmt"
	"time"
)

type OrderRepository interface {
	Store(order Order)
	FindById(id string) Order
}

const (
	MinUniquePremProductsForDiscount       int     = 3
	MaxUniqueProductsPerOrder              int     = 10
	DiscountValueIfThreeUniquePremProducts float64 = 0.1
)

type OrderError struct {
	Err error  // Error to bubble up and log somewhere
	Msg string // Error to show to client
}

func (e *OrderError) Error() string {
	return e.Err.Error()
}

type Order struct {
	id             string
	products       []Product
	productToCount map[Product]int
	dispatchDate   string
	status         OrderStatus
}

func NewOrder(id string) Order {
	return Order{
		id:             id,
		products:       make([]Product, 0),
		productToCount: make(map[Product]int),
	}
}

func (order *Order) ID() string {
	return order.id
}

func (order *Order) Value() float64 {
	sum := 0.0
	orderedProducts := order.products

	for _, product := range orderedProducts {
		sum += product.price
	}

	for _, product := range orderedProducts {
		if product.category == Premium && order.productToCount[product] >= MinUniquePremProductsForDiscount {
			sum *= (1 - DiscountValueIfThreeUniquePremProducts)
			break
		}
	}

	return sum
}

func (order *Order) Add(product Product) error {
	if !product.IsAvailable() {
		return &OrderError{Msg: fmt.Sprintf("product: %s cannot be added to the order as it is not available", product.name)}
	}

	if order.productToCount[product] > MaxUniqueProductsPerOrder {
		return &OrderError{Msg: fmt.Sprintf("product: %s cannot be added to the order as it exceeds the maximum allowed quantity per order, i.e., %d", product.name, MaxUniqueProductsPerOrder)}
	}

	order.productToCount[product] += 1
	order.products = append(order.products, product)

	return nil
}

func (order *Order) ProductQuantity() int {
	return len(order.products)
}

func (order *Order) Products() []Product {
	return order.products
}

func (order *Order) ProductToCount() map[Product]int {
	return order.productToCount
}

func (order *Order) SetDispatchDate(dateString string) error {
	if order.status != OrderDispatched {
		return &OrderError{Msg: "cannot set the dispatch date as order is not yet disptached"}
	}
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return &OrderError{Err: err, Msg: "invalid dispatch date format. please provide the correct date"}
	}
	if !date.After(time.Now()) {
		return &OrderError{Err: err, Msg: "dispatch date must be after the current date"}
	}
	order.dispatchDate = dateString
	return nil
}

func (order *Order) GetDispatchDate() string {
	return order.dispatchDate
}

func (order *Order) SetOrderStatus(status OrderStatus) {
	order.status = status
}

func (order *Order) GetOrderStatus() OrderStatus {
	return order.status
}
