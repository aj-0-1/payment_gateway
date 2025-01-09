package gateway

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type HTTPHandler struct {
	gateway PaymentGateway
}

func NewHTTPHandler(gateway PaymentGateway) *HTTPHandler {
	return &HTTPHandler{gateway: gateway}
}

var validate = validator.New()

func (h *HTTPHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var req ProcessPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payment, err := h.gateway.ProcessPayment(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

func (h *HTTPHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	payment, err := h.gateway.GetPayment(r.Context(), vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}
