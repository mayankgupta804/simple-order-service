package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type OrderRepository interface {
	Store(order Order) error
	FindById(id string) Order
	GetAll() []Order // Ideally this API should return paginated results, but for the sake of simplicity, we will be returning all results, assuming results are not many
}

const (
	MinUniquePremProductsForDiscount       int     = 3
	MaxUniqueProductsPerOrder              int     = 10
	DiscountValueIfThreeUniquePremProducts float64 = 0.1
)

type OrderError struct {
	Err error
}

func (e OrderError) Error() string {
	return e.Err.Error()
}

type Order struct {
	id             string
	products       []Product
	productToCount map[string]int
	dispatchDate   string
	status         OrderStatus
}

func NewOrder(id string) Order {
	return Order{
		id:             id,
		products:       make([]Product, 0),
		productToCount: make(map[string]int),
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
		if product.category == Premium && order.productToCount[product.id] >= MinUniquePremProductsForDiscount {
			sum *= (1 - DiscountValueIfThreeUniquePremProducts)
			break
		}
	}

	return sum
}

func (order *Order) Add(product Product) error {
	if !product.IsAvailable() {
		return &OrderError{Err: fmt.Errorf("product: %s cannot be added to the order as it is not available", product.name)}
	}

	if order.productToCount[product.id] > MaxUniqueProductsPerOrder {
		return &OrderError{Err: fmt.Errorf("product: %s cannot be added to the order as it exceeds the maximum allowed quantity per order, i.e., %d", product.name, MaxUniqueProductsPerOrder)}
	}

	order.productToCount[product.id] += 1
	order.products = append(order.products, product)

	return nil
}

func (order *Order) ProductQuantity() int {
	return len(order.products)
}

func (order *Order) Products() []Product {
	return order.products
}

func (order *Order) ProductToCount() map[string]int {
	return order.productToCount
}

func (order *Order) SetDispatchDate(dateString string) error {
	if order.status != OrderDispatched {
		return &OrderError{Err: errors.New("cannot set the dispatch date as order is not yet dispatched")}
	}
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return &OrderError{Err: errors.New("invalid dispatch date format. please provide the correct date")}
	}
	if !date.After(time.Now()) {
		return &OrderError{Err: errors.New("dispatch date must be after the current date")}
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

func (order *Order) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(struct {
		Id             string         `json:"id"`
		Products       []Product      `json:"products"`
		ProductToCount map[string]int `json:"product_to_count"`
		DispatchDate   string         `json:"dispatch_date"`
		Status         OrderStatus    `json:"status"`
	}{
		Id:             order.id,
		Products:       order.products,
		ProductToCount: order.productToCount,
		DispatchDate:   order.dispatchDate,
		Status:         order.status,
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (order *Order) UnmarshalJSON(data []byte) error {
	type ord struct {
		Id             string         `json:"id"`
		Products       []Product      `json:"products"`
		ProductToCount map[string]int `json:"product_to_count"`
		DispatchDate   string         `json:"dispatch_date"`
		Status         OrderStatus    `json:"status"`
	}
	o := &ord{}
	if err := json.Unmarshal(data, o); err != nil {
		return err
	}
	order.id = o.Id
	order.dispatchDate = o.DispatchDate
	order.productToCount = o.ProductToCount
	order.products = o.Products
	order.status = o.Status
	return nil
}
