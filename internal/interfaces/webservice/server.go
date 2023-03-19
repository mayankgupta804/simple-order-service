package webservice

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(orderInteractor OrderInteractor, productInteractor ProductInteractor) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		str := `{"status": "OK"}`
		io.WriteString(w, str)
	}).Methods(http.MethodGet)
	router.Handle("/orders", NewGetAllOrdersHandler(orderInteractor)).Methods(http.MethodGet)
	router.Handle("/orders/{id}", NewUpdateOrderHandler(orderInteractor)).Methods(http.MethodPut)
	router.Handle("/orders/{id}/products", NewAddProductToOrderHandler(orderInteractor)).Methods(http.MethodPost)
	router.Handle("/orders/{id}/products", NewGetAllOrderedProductsHandler(orderInteractor)).Methods(http.MethodGet)
	router.Handle("/products", NewGetAllProductsHandler(productInteractor)).Methods(http.MethodGet)
	router.Handle("/products/{id}", NewGetProductDetailsHandler(productInteractor)).Methods(http.MethodGet)
	return router
}

func StartServer(router *mux.Router) error {
	log.Printf("Starting web server on port: %v\n", "8080")
	if err := http.ListenAndServe(":"+"8080", router); err != nil {
		return err
	}
	return nil
}
