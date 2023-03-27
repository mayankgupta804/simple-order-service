package webservice

import (
	"encoding/json"
	"log"
	"net/http"
	"simple-order-service/internal/domain"
	"simple-order-service/internal/serializer"
	"simple-order-service/internal/usecases"
	"strings"

	"github.com/gorilla/mux"
)

type OrderInteractor interface {
	Products(orderId string) ([]usecases.Product, error)
	Add(orderId, productId string) error
	GetDetails(orderId string) (usecases.Order, error)
	GetAll() []usecases.Order
	UpdateDispatchDate(orderId, date string) error
	UpdateOrderStatus(orderId string, status domain.OrderStatus) error
}

type UpdateOrderHandler struct {
	orderInteractor OrderInteractor
}

type GetAllOrdersHandler struct {
	orderInteractor OrderInteractor
}

type GetAllOrderedProductsHandler struct {
	orderInteractor OrderInteractor
}

type GetOrderDetailsHandler struct {
	orderInteractor OrderInteractor
}

type AddProductToOrderHandler struct {
	orderInteractor OrderInteractor
}

func NewUpdateOrderHandler(orderInteractor OrderInteractor) UpdateOrderHandler {
	return UpdateOrderHandler{orderInteractor: orderInteractor}
}

func NewGetAllOrdersHandler(orderInteractor OrderInteractor) GetAllOrdersHandler {
	return GetAllOrdersHandler{orderInteractor: orderInteractor}
}

func NewGetAllOrderedProductsHandler(orderInteractor OrderInteractor) GetAllOrderedProductsHandler {
	return GetAllOrderedProductsHandler{orderInteractor: orderInteractor}
}

func NewGetOrderDetailsHandler(orderInteractor OrderInteractor) GetOrderDetailsHandler {
	return GetOrderDetailsHandler{orderInteractor: orderInteractor}
}

func NewAddProductToOrderHandler(orderInteractor OrderInteractor) AddProductToOrderHandler {
	return AddProductToOrderHandler{orderInteractor: orderInteractor}
}

func (handler GetOrderDetailsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	orderID := vars["id"]

	orderDetails, err := handler.orderInteractor.GetDetails(orderID)
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

	responseJSON, err := json.Marshal(orderDetails)
	if err != nil {
		log.Println(err.Error())
		failureResponse := serializer.Response{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(failureResponse.ToJSON())
		return
	}

	w.Write(responseJSON)
}

func (handler GetAllOrderedProductsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	orderID := vars["id"]

	orders, err := handler.orderInteractor.Products(orderID)
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

	responseJSON, err := json.Marshal(orders)
	if err != nil {
		log.Println(err.Error())
		failureResponse := serializer.Response{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(failureResponse.ToJSON())
		return
	}

	w.Write(responseJSON)
}

func (handler UpdateOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	orderID := vars["id"]

	decoder := json.NewDecoder(r.Body)

	var req serializer.UpdateOrderRequest

	if err := decoder.Decode(&req); err != nil {
		log.Println(err.Error())
		failureResponse := serializer.Response{
			Status:  "error",
			Message: "unable to parse JSON data",
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(failureResponse.ToJSON())
		return
	}

	res := serializer.Response{}
	res.Meta = &serializer.Meta{}
	res.Meta.Errors = []serializer.ErrorInfo{}

	var errorInfo serializer.ErrorInfo
	errCount := 0
	validUpdates := 0
	if len(strings.TrimSpace(req.DispatchDate)) > 0 {
		if err := handler.orderInteractor.UpdateDispatchDate(orderID, req.DispatchDate); err != nil {
			log.Println(err.Error())
			errorInfo = serializer.ErrorInfo{
				Detail: err.Error(),
			}
			res.Meta.Errors = append(res.Meta.Errors, errorInfo)
			errCount += 1
		} else {
			validUpdates += 1
		}
	}

	if len(strings.TrimSpace(req.OrderStatus)) > 0 {
		if err := handler.orderInteractor.UpdateOrderStatus(orderID, domain.OrderStatus(req.OrderStatus)); err != nil {
			log.Println(err.Error())
			errorInfo = serializer.ErrorInfo{
				Detail: err.Error(),
			}
			res.Meta.Errors = append(res.Meta.Errors, errorInfo)
			errCount += 1
		} else {
			validUpdates += 1
		}
	}

	if validUpdates == 0 {
		failureResponse := serializer.Response{
			Status:  "failure",
			Message: "update operations failed. more details can be found in the 'errors' section",
			Meta:    res.Meta,
		}
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(failureResponse.ToJSON())
		return
	} else if validUpdates <= errCount {
		partialSuccessResponse := serializer.Response{
			Status:  "partial success",
			Message: "one of the update operations failed. more details can be found in the 'errors' section",
			Meta:    res.Meta,
		}
		w.Header().Add("Content-type", "application/json; ext=partialsuccess")
		w.WriteHeader(http.StatusOK)
		w.Write(partialSuccessResponse.ToJSON())
		return
	}

	successResponse := serializer.Response{
		Status:  "success",
		Message: "order updated",
	}

	w.WriteHeader(http.StatusOK)
	w.Write(successResponse.ToJSON())
}

func (handler AddProductToOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	vars := mux.Vars(r)
	orderID := vars["id"]

	var req serializer.AddProductToOrderRequest
	if err := decoder.Decode(&req); err != nil {
		log.Println(err.Error())
		failureResponse := serializer.Response{
			Status:  "error",
			Message: "unable to parse JSON data",
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(failureResponse.ToJSON())
		return
	}

	if err := handler.orderInteractor.Add(orderID, req.ProductID); err != nil {
		log.Println(err.Error())
		failureResponse := serializer.Response{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(failureResponse.ToJSON())
		return
	}

	successResponse := serializer.Response{
		Status:  "success",
		Message: "product added to order",
	}

	w.WriteHeader(http.StatusOK)
	w.Write(successResponse.ToJSON())
}

func (handler GetAllOrdersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	orders := handler.orderInteractor.GetAll()

	responseJSON, err := json.Marshal(orders)

	if err != nil {
		log.Println(err.Error())
		failureResponse := serializer.Response{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(failureResponse.ToJSON())
		return
	}

	w.Write(responseJSON)
}
