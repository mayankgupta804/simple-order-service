package domain_test

import (
	"simple-order-service/internal/domain"
	"testing"
)

func TestDecreaseProductStock(t *testing.T) {
	product := domain.NewProduct("1", "nike shoes", 100.0, 5, domain.Premium)

	if product.SKU() != 5 {
		t.Errorf("Got: %v, Want: %v", product.SKU(), 5)
	}

	product.DecreaseStockBy(1)

	if product.SKU() != 4 {
		t.Errorf("Got: %v, Want: %v", product.SKU(), 4)
	}
}

func TestDecreaseProductStockIfStockIsZero(t *testing.T) {
	product := domain.NewProduct("1", "nike shoes", 100.0, 0, domain.Premium)

	product.DecreaseStockBy(1)
	got := product.SKU()
	want := 0

	if got != want {
		t.Errorf("Got: %v, Want: %v", got, want)
	}
}

func TestDecreaseProductIsAvailable_True(t *testing.T) {
	product := domain.NewProduct("1", "nike shoes", 100.0, 1, domain.Premium)

	got := product.IsAvailable()
	want := true

	if got != want {
		t.Errorf("Got: %v, Want: %v", got, want)
	}
}

func TestDecreaseProductIsAvailable_False(t *testing.T) {
	product := domain.NewProduct("1", "nike shoes", 100.0, 0, domain.Premium)

	got := product.IsAvailable()
	want := false

	if got != want {
		t.Errorf("Got: %v, Want: %v", got, want)
	}
}
