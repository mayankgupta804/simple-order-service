package domain

import "encoding/json"

type ProductRepository interface {
	Store(product Product) error
	FindById(id string) Product
	GetAll() []Product // Ideally this API should return paginated results, but for the sake of simplicity, we will be returning all results, assuming results are not many
}

type Product struct {
	id       string
	name     string
	price    float64
	sku      int
	category ProductCategory
}

func NewProduct(id string, name string, price float64, sku int, category ProductCategory) Product {
	return Product{
		id:       id,
		name:     name,
		price:    price,
		category: category,
		sku:      sku,
	}
}

func (product *Product) Name() string {
	return product.name
}

func (product *Product) ID() string {
	return product.id
}

func (product *Product) Price() float64 {
	return product.price
}

func (product *Product) Category() ProductCategory {
	return product.category
}

func (product *Product) DecreaseStockBy(decreaseBy int) {
	currentStock := product.sku
	currentStock -= decreaseBy
	if currentStock <= 0 {
		product.sku = 0
		return
	}
	product.sku = currentStock
}

func (product *Product) IsAvailable() bool {
	return product.sku > 0
}

func (product *Product) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(struct {
		Id       string
		Name     string
		Price    float64
		Sku      int
		Category ProductCategory
	}{
		Id:       product.id,
		Name:     product.name,
		Price:    product.price,
		Sku:      product.sku,
		Category: product.category,
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}
