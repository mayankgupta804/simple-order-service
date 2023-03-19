package main

import (
	"log"
	"os"
	"simple-order-service/internal/domain"
	"simple-order-service/internal/interfaces/repository"
	"simple-order-service/internal/interfaces/webservice"
	"simple-order-service/internal/usecases"
	"simple-order-service/pkg/database"

	"github.com/urfave/cli"
)

func main() {

	clientApp := cli.NewApp()
	clientApp.Name = "Simple Order Service"
	clientApp.Version = "0.0.1"
	clientApp.Commands = []cli.Command{
		{
			Name:        "start:webserver",
			Description: "Start Webserver",
			Action: func(c *cli.Context) {
				StartWebServer()
			},
		},
		{
			Name:        "seed:db:products",
			Description: "Seed products to DB",
			Action: func(c *cli.Context) {
				SeedProductsInDB()
			},
		},
	}

	if err := clientApp.Run(os.Args); err != nil {
		panic(err)
	}
}

func StartWebServer() {
	db, err := database.NewInstance("shop.db")
	if err != nil {
		log.Fatal(err)
	}
	var ordersRepo domain.OrderRepository = repository.NewOrdersRepo(db)
	var productsRepo domain.ProductRepository = repository.NewProductsRepo(db)

	var orderInteractor webservice.OrderInteractor = usecases.NewOrderInteractor(ordersRepo, productsRepo)
	var productInteractor webservice.ProductInteractor = usecases.NewProductInteractor(productsRepo)

	router := webservice.SetupRoutes(orderInteractor, productInteractor)

	if err = webservice.StartServer(router); err != nil {
		log.Fatal(err)
	}
}

func SeedProductsInDB() {
	db, err := database.NewInstance("shop.db")
	if err != nil {
		log.Fatal(err)
	}

	var productsRepo domain.ProductRepository = repository.NewProductsRepo(db)
	var productInteractor webservice.ProductInteractor = usecases.NewProductInteractor(productsRepo)

	product1 := domain.NewProduct("1", "sneakers", 12.0, 11, domain.Premium)
	product2 := domain.NewProduct("2", "shirt", 10.0, 3, domain.Premium)
	product3 := domain.NewProduct("3", "trousers", 20.0, 5, domain.Premium)
	product4 := domain.NewProduct("4", "tie", 10.0, 12, domain.Budget)

	data1, _ := product1.MarshalJSON()
	data2, _ := product2.MarshalJSON()
	data3, _ := product3.MarshalJSON()
	data4, _ := product4.MarshalJSON()

	db.Put([]byte("products"), []byte(product1.ID()), data1)
	db.Put([]byte("products"), []byte(product2.ID()), data2)
	db.Put([]byte("products"), []byte(product3.ID()), data3)
	db.Put([]byte("products"), []byte(product4.ID()), data4)

	log.Println(productInteractor.GetAll())
}
