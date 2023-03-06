package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type storage interface {
	LockResource(context.Context, string, string, time.Duration) (bool, error)
	UnlockResource(context.Context, string, string) error
}

type Handler struct {
	storage storage
}

func NewHandler(s storage) Handler {
	return Handler{
		storage: s,
	}
}

type LockResponse struct {
	RandomValue string `json:"random_value,omitempty"`
	Succeed     bool   `json:"succeed"`
}

func (h *Handler) Lock(w http.ResponseWriter, r *http.Request) {
	resource, ok := mux.Vars(r)["resource"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("resource is absent in request")
		return
	}

	durationStr := r.URL.Query().Get("duration")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("duration is incorrect or absent in request")
		return
	}

	randomValue := uuid.New().String()
	succeed, err := h.storage.LockResource(r.Context(), resource, randomValue, duration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to lock resource: %s", err.Error())
		return
	}

	var response LockResponse
	if succeed {
		response = LockResponse{
			Succeed:     succeed,
			RandomValue: randomValue,
		}
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to marshal response: %s", err.Error())
		return
	}

	_, _ = w.Write(responseData)
}

func (h *Handler) Unlock(w http.ResponseWriter, r *http.Request) {
	resource, ok := mux.Vars(r)["resource"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("resource is absent in request")
		return
	}

	randomValue := r.URL.Query().Get("randomValue")
	if randomValue == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("random value is absent in request")
		return
	}

	if err := h.storage.UnlockResource(r.Context(), resource, randomValue); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to unlock resource: %s", err.Error())
		return
	}
}
