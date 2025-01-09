package main

import (
	"log"
	"net/http"
	"payment_gateway/internal/bank"
	"payment_gateway/internal/gateway"
	"payment_gateway/internal/store"

	"github.com/gorilla/mux"
)

func main() {
	store := store.NewMemoryStore()
	bankService := bank.NewSimulator()
	gatewayService := gateway.NewService(store, bankService)

	handler := gateway.NewHTTPHandler(gatewayService)

	router := mux.NewRouter()
	router.HandleFunc("/payments", handler.ProcessPayment).Methods("POST")
	router.HandleFunc("/payments/{id}", handler.GetPayment).Methods("GET")

	log.Printf("Starting payment gateway server")
	log.Fatal(http.ListenAndServe(":8080", router))
}
