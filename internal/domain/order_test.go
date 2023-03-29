package domain_test

import (
	"errors"
	"simple-order-service/internal/domain"
	"testing"
	"time"
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

func TestCreateOrderWithThreeUniqueProducts(t *testing.T) {
	orderID := "123"

	order := domain.NewOrder(orderID)

	product1 := domain.NewProduct("1", "nike shoes", 100.0, 3, domain.Premium)
	product2 := domain.NewProduct("2", "adidas shoes", 50.0, 2, domain.Premium)
	product3 := domain.NewProduct("3", "puma shoes", 10.0, 2, domain.Premium)

	order.Add(product1)
	order.Add(product2)
	order.Add(product3)

	totalProductPrice := (product1.Price() + product2.Price() + product3.Price())

	expectedOrderValue := totalProductPrice * float64(1-domain.DiscountValueIfThreeUniquePremProducts)
	actualOrderValue := order.Value()

	if actualOrderValue != expectedOrderValue {
		t.Errorf("Expected order value: %v. Got order value: %v", expectedOrderValue, actualOrderValue)
	}
}

func TestCreateOrderWithDispatchDate_DispatchError(t *testing.T) {
	orderID := "123"
	order := domain.NewOrder(orderID)
	product1 := domain.NewProduct("1", "nike shoes", 100.0, 3, domain.Premium)
	order.Add(product1)

	got := order.SetDispatchDate(time.Now().Add(time.Duration(time.Now().Day())).Format(time.DateOnly))
	want := domain.OrderError{Err: domain.ErrInvalidDispatchDateWithOrderNotDispatched}

	if got.Error() != want.Error() {
		t.Errorf("Got: %v, Want: %v", got, want)
	}
}

func TestCreateOrderWithDispatchDate_InvalidDateFormat(t *testing.T) {
	orderID := "123"
	order := domain.NewOrder(orderID)
	product1 := domain.NewProduct("1", "nike shoes", 100.0, 3, domain.Premium)
	order.Add(product1)
	order.SetOrderStatus(domain.OrderDispatched)

	got := order.SetDispatchDate("2023-31-12")
	want := domain.OrderError{Err: domain.ErrInvalidDispatchDateFormat}
	if got.Error() != want.Error() {
		t.Errorf("Got: %v, Want: %v", got, want)
	}
}

func TestCreateOrderWithDispatchDate_InvalidDispatchDate(t *testing.T) {
	orderID := "123"
	order := domain.NewOrder(orderID)
	product1 := domain.NewProduct("1", "nike shoes", 100.0, 3, domain.Premium)
	order.Add(product1)
	order.SetOrderStatus(domain.OrderDispatched)

	got := order.SetDispatchDate(time.Now().Add(-time.Duration(time.Now().Day())).Format(time.DateOnly))
	want := domain.OrderError{Err: domain.ErrInvalidDispatchDate}
	if got.Error() != want.Error() {
		t.Errorf("Got: %v, Want: %v", got, want)
	}
}

func TestCreateOrderWithDispatchDate_Success(t *testing.T) {
	orderID := "123"
	order := domain.NewOrder(orderID)
	product1 := domain.NewProduct("1", "nike shoes", 100.0, 3, domain.Premium)
	order.Add(product1)
	order.SetOrderStatus(domain.OrderDispatched)

	got := order.SetDispatchDate(time.Now().Add(48 * time.Hour).Format(time.DateOnly))
	var want error = nil
	if got != want {
		t.Errorf("Got: %v, Want: %v", got, want)
	}
}

func TestAddProductToOrder_ProductUnavailable(t *testing.T) {
	orderID := "123"
	order := domain.NewOrder(orderID)
	product1 := domain.NewProduct("1", "nike shoes", 100.0, 0, domain.Premium)
	got := order.Add(product1)
	want := domain.OrderError{Err: domain.ErrUnavailableProduct(product1.Name())}
	var orderErr *domain.OrderError

	if !errors.As(got, &orderErr) {
		t.Errorf("Got: %v, Want: %v", got, orderErr)
	}

	if got.Error() != want.Error() {
		t.Errorf("Got: %v, Want: %v", got, want)
	}
}

func TestAddProductToOrder_MaxAllowedQuantityError(t *testing.T) {
	orderID := "123"
	order := domain.NewOrder(orderID)
	product1 := domain.NewProduct("1", "nike shoes", 100.0, 11, domain.Premium)
	for i := 1; i <= 10; i++ {
		order.Add(product1)
	}
	got := order.Add(product1)
	want := domain.OrderError{Err: domain.ErrMaxAllowedQuantity(product1.Name())}
	var orderErr *domain.OrderError

	if !errors.As(got, &orderErr) {
		t.Errorf("Got: %v, Want: %v", got, orderErr)
	}

	if got.Error() != want.Error() {
		t.Errorf("Got: %v, Want: %v", got, want)
	}
}
