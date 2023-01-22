package app

import (
	"banking/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CustomerHandlers struct {
	service service.CustomerService
}

func (ch *CustomerHandlers) getAllCustomers(w http.ResponseWriter, r *http.Request) {

	status := r.URL.Query().Get("status")

	customers, err := ch.service.GetAllCustomer(status)

	if err != nil {
		w.WriteHeader(err.Code)
		errorJson, _ := json.Marshal(err.AsMessage())
		w.Write([]byte(errorJson))
	} else {
		w.WriteHeader(err.Code)
		customer, _ := json.Marshal(customers)
		w.Write([]byte(customer))
	}
}

func (ch *CustomerHandlers) getCustomer(w http.ResponseWriter, r *http.Request) {
	Id := chi.URLParam(r, "customer_id")

	customer, err := ch.service.GetCustomer(Id)
	if err != nil {
		w.WriteHeader(err.Code)
		w.Write([]byte(err.Message))
	} else {
		w.WriteHeader(http.StatusOK)
		customerJSON, _ := json.Marshal(customer)
		w.Write(customerJSON)
	}
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
