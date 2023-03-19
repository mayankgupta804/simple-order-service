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

func NewOrderInteractor(orderRepo domain.OrderRepository, productRepo domain.ProductRepository) *OrderInteractor {
	return &OrderInteractor{orderRepository: orderRepo, productRepository: productRepo}
}

func (interactor *OrderInteractor) Products(orderId string) ([]Product, error) {
	order := interactor.orderRepository.FindById(orderId)
	orderedProducts := order.Products()
	if len(orderedProducts) == 0 {
		return nil, errors.New("order does not exist. no products found in the order")
	}

	products := make([]Product, len(orderedProducts))
	for idx, product := range orderedProducts {
		products[idx] = Product{ID: product.ID(), Name: product.Name(), Category: string(product.Category()), Price: product.Price()}
	}
	return products, nil
}

func (interactor *OrderInteractor) Add(orderId, productId string) error {
	product := interactor.productRepository.FindById(productId)
	order := interactor.orderRepository.FindById(orderId)
	if order.ID() == "" {
		order = domain.NewOrder(orderId)
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
	return nil
}

func (interactor *OrderInteractor) UpdateOrderStatus(orderId string, status domain.OrderStatus) error {
	// TODO: Also validate status
	// TODO: Always check previous status; status can only move forwards, i.e., placed -> dispatched -> completed
	order := interactor.orderRepository.FindById(orderId)

	if order.ID() == "" {
		return errors.New("cannot update order status for a non-existent order")
	}

	if status == domain.OrderPlaced {
		for productId, count := range order.ProductToCount() {
			product := interactor.productRepository.FindById(productId)
			product.DecreaseStockBy(count)
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
