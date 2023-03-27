package usecases

import (
	"errors"
	"fmt"
	"simple-order-service/internal/domain"
)

type OrderInteractor struct {
	orderRepository   domain.OrderRepository
	productRepository domain.ProductRepository
}

type Order struct {
	ID            string    `json:"id"`
	TotalQuantity int       `json:"total_quantity"`
	Products      []Product `json:"products"`
	DispatchDate  string    `json:"dispatch_date,omitempty"`
	Status        string    `json:"status,omitempty"`
	Value         float64   `json:"value,omitempty"`
}

func NewOrderInteractor(orderRepo domain.OrderRepository, productRepo domain.ProductRepository) *OrderInteractor {
	return &OrderInteractor{orderRepository: orderRepo, productRepository: productRepo}
}

// TODO: Refactor after writing the deduplication of products in the domain.Orders.Add()
func (interactor *OrderInteractor) Products(orderId string) ([]Product, error) {
	order := interactor.orderRepository.FindById(orderId)
	orderedProducts := order.Products()
	if order.ProductQuantity() == 0 {
		return nil, errors.New("order does not exist. no products found in the order")
	}

	products := make([]Product, 0)
	seenProducts := make(map[string]bool)
	productToCount := make(map[Product]int)

	for _, product := range orderedProducts {
		prod := Product{ID: product.ID(), Name: product.Name(), Category: string(product.Category()), Price: product.Price()}
		if seenProducts[product.ID()] {
			productToCount[prod] += 1
			continue
		}
		seenProducts[product.ID()] = true
		productToCount[prod] += 1
	}

	for product, count := range productToCount {
		if count > 0 {
			product.Quantity = count
			products = append(products, product)
		}
	}
	return products, nil
}

func (interactor *OrderInteractor) Add(orderId, productId string) error {
	product := interactor.productRepository.FindById(productId)
	order := interactor.orderRepository.FindById(orderId)
	if order.ID() == "" {
		order = domain.NewOrder(orderId)
	}

	orderStatus := order.GetOrderStatus()
	if orderStatus == domain.OrderCompleted || orderStatus == domain.OrderCancelled || orderStatus == domain.OrderDispatched {
		return fmt.Errorf("order has already been %s", orderStatus)
	}

	if domainErr := order.Add(product); domainErr != nil {
		message := "Could not add item #%s "
		message += "to order #%s "
		message += "because a business rule was violated: '%s'"
		err := fmt.Errorf(message,
			product.ID(),
			order.ID(),
			domainErr.Error())
		return err
	}
	interactor.orderRepository.Store(order)
	interactor.UpdateOrderStatus(orderId, domain.OrderPlaced)
	return nil
}

// TODO: Always check previous status; status can only move forwards, i.e., placed -> dispatched or cancelled -> completed
func (interactor *OrderInteractor) UpdateOrderStatus(orderId string, status domain.OrderStatus) error {
	orderStatusMap := map[domain.OrderStatus]bool{
		domain.OrderDispatched: true,
		domain.OrderPlaced:     true,
		domain.OrderCompleted:  true,
		domain.OrderCancelled:  true,
	}
	_, ok := orderStatusMap[status]
	if !ok {
		return errors.New(`invalid order status. the different order status values are: 
							'placed', 'dispatched', 'cancelled' and 'completed'`)
	}

	order := interactor.orderRepository.FindById(orderId)
	orderStatus := order.GetOrderStatus()
	if orderStatus == domain.OrderCancelled || orderStatus == domain.OrderDispatched {
		return fmt.Errorf("cannot update order status as order has been %s", orderStatus)
	}

	if order.ID() == "" {
		return errors.New("cannot update order status for a non-existent order")
	}

	if status == domain.OrderPlaced {
		for productId := range order.ProductToCount() {
			product := interactor.productRepository.FindById(productId)
			product.DecreaseStockBy(1)
			interactor.productRepository.Store(product)
		}
	}

	order.SetOrderStatus(status)
	interactor.orderRepository.Store(order)
	return nil
}

func (interactor *OrderInteractor) UpdateDispatchDate(orderId, date string) error {
	var message string
	order := interactor.orderRepository.FindById(orderId)
	if order.ID() == "" {
		return errors.New("cannot update dispatch date for a non-existent order")
	}

	orderStatus := order.GetOrderStatus()
	if orderStatus == domain.OrderCompleted || orderStatus == domain.OrderCancelled {
		return fmt.Errorf("cannot update dispatch date as order has been %s", orderStatus)
	}

	if domainErr := order.SetDispatchDate(date); domainErr != nil {
		message = "Could not update dispatch date: #%s "
		message += "of order #%s "
		message += "because a business rule was violated: '%s'"
		err := fmt.Errorf(message,
			date,
			order.ID(),
			domainErr.Error())
		return err
	}
	interactor.orderRepository.Store(order)
	return nil
}

func (interactor *OrderInteractor) GetDetails(orderId string) (Order, error) {
	domainOrder := interactor.orderRepository.FindById(orderId)
	if domainOrder.ID() == "" {
		return Order{}, errors.New("order does not exist")
	}
	order := Order{
		ID:            domainOrder.ID(),
		TotalQuantity: domainOrder.ProductQuantity(),
		DispatchDate:  domainOrder.GetDispatchDate(),
		Status:        string(domainOrder.GetOrderStatus()),
		Value:         domainOrder.Value(),
		Products:      getDeduplicatedProductsWithCount(domainOrder.Products()),
	}
	return order, nil
}

func (interactor *OrderInteractor) GetAll() []Order {
	ordersFromDb := interactor.orderRepository.GetAll()
	if len(ordersFromDb) == 0 {
		return []Order{}
	}
	orders := make([]Order, len(ordersFromDb))
	for idx, order := range ordersFromDb {
		orders[idx] = Order{
			ID:            order.ID(),
			TotalQuantity: order.ProductQuantity(),
			DispatchDate:  order.GetDispatchDate(),
			Status:        string(order.GetOrderStatus()),
			Value:         order.Value(),
			Products:      getDeduplicatedProductsWithCount(order.Products()),
		}
	}
	return orders
}

func getDeduplicatedProductsWithCount(products []domain.Product) []Product {
	productToCount := make(map[domain.Product]int)
	for _, product := range products {
		productToCount[product] += 1
	}
	deduplicatedProducts := make([]Product, 0)
	for product, count := range productToCount {
		p := Product{
			ID:       product.ID(),
			Name:     product.Name(),
			Category: string(product.Category()),
			Price:    product.Price(),
			Quantity: count,
		}
		deduplicatedProducts = append(deduplicatedProducts, p)
	}
	return deduplicatedProducts
}
