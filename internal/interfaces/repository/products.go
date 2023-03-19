package repository

import (
	"encoding/json"
	"simple-order-service/internal/domain"
	"simple-order-service/pkg/database"
)

const ProductsSchema = "products"

type productsRepo struct {
	dbClient database.DB
}

func NewProductsRepo(db *database.DB) productsRepo {
	return productsRepo{dbClient: *db}
}

func (prodRepo productsRepo) Store(product domain.Product) error {
	data, err := json.Marshal(&product)
	if err != nil {
		return err
	}
	return prodRepo.dbClient.Put([]byte(ProductsSchema), []byte(product.ID()), data)
}

func (prodRepo productsRepo) FindById(id string) domain.Product {
	product := &domain.Product{}
	data := prodRepo.dbClient.Get([]byte(ProductsSchema), []byte(id))
	if data == nil {
		return *product
	}
	product.UnmarshalJSON(data)
	return *product
}

func (prodRepo productsRepo) GetAll() []domain.Product {
	data := prodRepo.dbClient.GetAll([]byte(ProductsSchema))
	if len(data) == 0 {
		return []domain.Product{}
	}
	products := make([]domain.Product, len(data))
	for idx, val := range data {
		product := &domain.Product{}
		product.UnmarshalJSON(val)
		products[idx] = *product
	}
	return products
}
