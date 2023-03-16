package usecases

import (
	"errors"
	"simple-order-service/internal/domain"
)

type Product struct {
	ID       string
	Name     string
	Category string
	Price    float64
}
type ProductInteractor struct {
	ProductRepository domain.ProductRepository
}

func (interactor *ProductInteractor) GetDetails(productID string) (domain.Product, error) {
	product := interactor.ProductRepository.FindById(productID)
	if product.ID() == "" {
		return domain.Product{}, errors.New("product does not exist")
	}
	return product, nil
}

func (interactor *ProductInteractor) GetAll() []Product {
	productsFromDb := interactor.ProductRepository.GetAll()
	products := make([]Product, len(productsFromDb))
	for _, product := range productsFromDb {
		products = append(products, Product{ID: product.ID(), Name: product.Name(), Category: string(product.Category()), Price: product.Price()})
	}
	return products
}
