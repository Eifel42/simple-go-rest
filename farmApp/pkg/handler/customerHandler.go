package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"farmApp/pkg/api"
	"farmApp/pkg/persistence"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// @Summary Get all customers
// @Description Get all customers
// @Tags customers
// @Produce json
// @Success 200 {array} api.Customer
// @Failure 500 {object} api.ErrorResponse
// @Router /customers [get]
func GetCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := persistence.GetCustomers()
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	encodeJSONResponse(w, customers)
}

// @Summary Get a customer by ID
// @Description Get a customer by ID
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} api.Customer
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /customers/{id} [get]
func GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	customer, err := persistence.GetCustomerByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Customer not found", http.StatusNotFound)
		} else {
			handleError(w, err, http.StatusInternalServerError)
		}
		return
	}
	encodeJSONResponse(w, customer)
}

// @Summary Add a new customer
// @Description Add a new customer
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body api.Customer true "Customer"
// @Success 201 {object} api.Customer
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /customers [post]
func AddCustomer(w http.ResponseWriter, r *http.Request) {
	var customer api.Customer

	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	id, err := persistence.AddCustomer(customer)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	customer.ID = &id
	w.WriteHeader(http.StatusCreated)
	encodeJSONResponse(w, customer)
}

// @Summary Update a customer
// @Description Update a customer
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body api.Customer true "Customer"
// @Success 200 {object} api.Customer
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /customers/{id} [put]
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer api.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}
	customer.ID = &id

	err = persistence.UpdateCustomer(*customer.ID, customer)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	encodeJSONResponse(w, customer)
}

// @Summary Delete a customer
// @Description Delete a customer
// @Tags customers
// @Param id path int true "Customer ID"
// @Success 204
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /customers/{id} [delete]
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	_, err = persistence.GetCustomerByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Customer not found", http.StatusNotFound)
		} else {
			handleError(w, err, http.StatusInternalServerError)
		}
		return
	}

	err = persistence.DeleteCustomer(id)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleError(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
}

func encodeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		handleError(w, err, http.StatusInternalServerError)
	}
}
