package domain_test

import (
	"simple-order-service/internal/domain"
	"testing"
)

func TestCreateOrderWithProducts(t *testing.T) {
	order := domain.NewOrder("123")
	if order.Products() == nil {
		t.Fail()
	}
	if order.ProductToCount() == nil {
		t.Fail()
	}
	if order.ID() != "123" {
		t.Fail()
	}

	product1 := domain.NewProduct("1", "nike shoes", 11.0, 3, domain.Premium)
	product2 := domain.NewProduct("2", "adidas shoes", 13.0, 2, domain.Premium)
	order.Add(product1)
	order.Add(product2)

	if order.Value() != 35.0 {
		t.Fail()
	}

	if order.ProductToCount()["1"] != 2 {
		t.Fail()
	}

	if len(order.Products()) != 3 {
		t.Fail()
	}
}
