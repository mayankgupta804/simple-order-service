package repository

import (
	"simple-order-service/internal/domain"
	"simple-order-service/pkg/database"
)

const OrdersSchema = "orders"

type ordersRepo struct {
	dbClient database.DB
}

func NewOrdersRepo(db *database.DB) ordersRepo {
	return ordersRepo{dbClient: *db}
}

func (ordRepo ordersRepo) Store(order domain.Order) error {
	data, err := order.MarshalJSON()
	if err != nil {
		return err
	}
	return ordRepo.dbClient.Put([]byte(OrdersSchema), []byte(order.ID()), data)
}

func (ordRepo ordersRepo) FindById(id string) domain.Order {
	data := ordRepo.dbClient.Get([]byte(OrdersSchema), []byte(id))
	order := &domain.Order{}
	order.UnmarshalJSON(data)
	return *order
}

func (ordRepo ordersRepo) GetAll() []domain.Order {
	data := ordRepo.dbClient.GetAll([]byte(OrdersSchema))
	if len(data) == 0 {
		return []domain.Order{}
	}
	orders := make([]domain.Order, len(data))
	for idx, val := range data {
		order := &domain.Order{}
		order.UnmarshalJSON(val)
		orders[idx] = *order
	}
	return orders
}
