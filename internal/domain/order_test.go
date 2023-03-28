package domain_test

import (
	"simple-order-service/internal/domain"
	"testing"
)

func TestCreateOrderWithProducts(t *testing.T) {
	orderID := "123"

	order := domain.NewOrder(orderID)
	if len(order.Products()) != 0 {
		t.Error("products in the order must be empty when a new order is created")
	}
	if order.ProductToCount() == nil {
		t.Error("the map of products in the order to the count must be initialized")
	}

	product1 := domain.NewProduct("1", "nike shoes", 11.0, 3, domain.Premium)
	product2 := domain.NewProduct("2", "adidas shoes", 13.0, 2, domain.Premium)

	order.Add(product1)
	order.Add(product2)

	expectedOrderValue := product1.Price() + product2.Price()
	actualOrderValue := order.Value()

	if actualOrderValue != expectedOrderValue {
		t.Errorf("Expected order value: %v. Got order value: %v", expectedOrderValue, actualOrderValue)
	}

	actualProduct1Count := order.ProductToCount()["1"]
	expectedProduct1Count := 1

	if actualProduct1Count != expectedProduct1Count {
		t.Errorf("Expected count of product1: %d. Got count of product1: %d", expectedProduct1Count, actualProduct1Count)
	}

	expectedTotalOrderedProducts := 2
	actualTotalOrderedProducts := len(order.Products())
	if expectedTotalOrderedProducts != actualTotalOrderedProducts {
		t.Errorf("Expected total products: %d. Got total products: %d", expectedProduct1Count, actualProduct1Count)
	}
}
