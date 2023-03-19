package usecases

import (
	"errors"
	"simple-order-service/internal/domain"
)

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

type ProductInteractor struct {
	productRepository domain.ProductRepository
}

func NewProductInteractor(productRepo domain.ProductRepository) *ProductInteractor {
	return &ProductInteractor{productRepository: productRepo}
}

func (interactor *ProductInteractor) GetDetails(productID string) (domain.Product, error) {
	product := interactor.productRepository.FindById(productID)
	if product.ID() == "" {
		return domain.Product{}, errors.New("product does not exist")
	}
	return product, nil
}

func (interactor *ProductInteractor) GetAll() []Product {
	productsFromDb := interactor.productRepository.GetAll()
	if len(productsFromDb) == 0 {
		return []Product{}
	}
	products := make([]Product, len(productsFromDb))
	for idx, product := range productsFromDb {
		products[idx] = Product{ID: product.ID(), Name: product.Name(), Category: string(product.Category()), Price: product.Price()}
	}
	return products
}
