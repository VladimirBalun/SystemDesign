package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type storage interface {
	CreateOrder(Order) (int, error)
	RemoveOrder(int) error
	SubmitOrder(int) error
}

type saga interface {
	NotifyOrderCreated(int, []string) error
}

type HTTPHandler struct {
	storage storage
	saga    saga
}

func NewHTTPHandler(st storage, sg saga) HTTPHandler {
	return HTTPHandler{
		storage: st,
		saga:    sg,
	}
}

type CreateOrderRequest struct {
	Username string   `json:"username"`
	Goods    []string `json:"goods"`
}

func (h *HTTPHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request CreateOrderRequest
	if err = json.Unmarshal(body, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order := NewOrder(request.Username, request.Goods)
	orderID, err := h.storage.CreateOrder(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = h.saga.NotifyOrderCreated(orderID, order.Goods()); err != nil {
		if err = h.storage.RemoveOrder(orderID); err != nil {
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
