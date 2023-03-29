package usecases_test

import (
	"simple-order-service/internal/domain"
	"simple-order-service/internal/usecases"
	"testing"
)

func TestListProductsInOrder(t *testing.T) {
	orderRepoMock := &domain.OrderRepositoryMock{
		GetAllFunc: func() []domain.Order {
			order := domain.NewOrder("1")
			product := domain.NewProduct("123", "nike shoes", 100.0, 5, domain.Premium)
			order.Add(product)
			return []domain.Order{order}
		},
	}
	productRepoMock := &domain.ProductRepositoryMock{}

	orderInteractor := usecases.NewOrderInteractor(orderRepoMock, productRepoMock)
	got := orderInteractor.GetAll()
	if len(got) != 1 {
		t.Error("number of orders must be equal to 1")
	}
	if got[0].ID != "1" {
		t.Error("id of order should be 1")
	}
}
