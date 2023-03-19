package webservice

import (
	"encoding/json"
	"log"
	"net/http"
	"simple-order-service/internal/serializer"
	"simple-order-service/internal/usecases"

	"github.com/gorilla/mux"
)

type ProductInteractor interface {
	GetDetails(productId string) (usecases.Product, error)
	GetAll() []usecases.Product
}

type GetProductDetailsHandler struct {
	productInteractor ProductInteractor
}

type GetAllProductsHandler struct {
	productInteractor ProductInteractor
}

func NewGetProductDetailsHandler(productInteractor ProductInteractor) GetProductDetailsHandler {
	return GetProductDetailsHandler{productInteractor: productInteractor}
}
func NewGetAllProductsHandler(productInteractor ProductInteractor) GetAllProductsHandler {
	return GetAllProductsHandler{productInteractor: productInteractor}
}

func (handler GetProductDetailsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	productID := vars["id"]

	productDetails, err := handler.productInteractor.GetDetails(productID)
	if err != nil {
		log.Println(err.Error())
		failureResponse := serializer.Response{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(failureResponse.ToJSON())
		return
	}

	responseJSON, _ := json.Marshal(productDetails)

	w.Write(responseJSON)
}

func (handler GetAllProductsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	products := handler.productInteractor.GetAll()

	responseJSON, _ := json.Marshal(products)

	w.Write(responseJSON)
}
