package usecases

import (
	"errors"
	"fmt"
	"simple-order-service/internal/domain"
)

type OrderInteractor struct {
	OrderRepository   domain.OrderRepository
	ProductRepository domain.ProductRepository
}

func (interactor *OrderInteractor) Products(orderId string) ([]Product, error) {
	order := interactor.OrderRepository.FindById(orderId)
	orderedProducts := order.Products()
	if orderedProducts == nil {
		return nil, errors.New("order does not exist. no products found in the order")
	}
	products := make([]Product, len(orderedProducts))
	for _, product := range orderedProducts {
		products = append(products, Product{ID: product.ID(), Name: product.Name(), Category: string(product.Category()), Price: product.Price()})
	}
	return products, nil
}

func (interactor *OrderInteractor) Add(orderId, productId string) error {
	product := interactor.ProductRepository.FindById(productId)
	order := interactor.OrderRepository.FindById(orderId)
	if order.ID() == "" {
		order = domain.NewOrder(orderId)
	}
	if domainErr := order.Add(product); domainErr != nil {
		message := "Could not add item #%s "
		message += "to order #%s "
		message += "because a business rule was violated: '%s'"
		err := fmt.Errorf(message,
			product.ID,
			order.ID,
			domainErr.Error())
		return err
	}
	interactor.OrderRepository.Store(order)
	return nil
}

func (interactor *OrderInteractor) UpdateOrderStatus(orderId string, status domain.OrderStatus) error {
	order := interactor.OrderRepository.FindById(orderId)

	if order.ID() == "" {
		return errors.New("cannot update order status for a non-existent order")
	}

	if status == domain.OrderPlaced {
		for product, count := range order.ProductToCount() {
			product.DecreaseStockBy(count)
			interactor.ProductRepository.Store(product)
		}
	}

	order.SetOrderStatus(status)
	interactor.OrderRepository.Store(order)
	return nil
}

func (interactor *OrderInteractor) UpdateDispatchDate(orderId, date string) error {
	var message string
	order := interactor.OrderRepository.FindById(orderId)
	if order.ID() == "" {
		return errors.New("cannot update dispatch date for a non-existent order")
	}
	if domainErr := order.SetDispatchDate(date); domainErr != nil {
		message = "Could not update dispatch date: #%s "
		message += "of order #%s "
		message += "because a business rule was violated: '%s'"
		err := fmt.Errorf(message,
			date,
			order.ID,
			domainErr.Error())
		return err
	}
	return nil
}
