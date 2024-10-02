package routeHandler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"farmApp/farmApp/pkg/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var db *sql.DB

// @Summary Get all pkg
// @Description Get all pkg
// @Tags pkg
// @Produce json
// @Success 200 {array} Customer
// @Router /pkg [get]
func getCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, role, email, phone, contacted FROM pkg")
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			handleError(w, err, http.StatusInternalServerError)
		}
	}(rows)

	var customers []model.Customer

	for rows.Next() {
		var customer model.Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Role, &customer.Email, &customer.Phone, &customer.Contacted); err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		customers = append(customers, customer)
	}

	encodeJSONResponse(w, customers)
}

// @Summary Get a customer by ID
// @Description Get a customer by ID
// @Tags pkg
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} Customer
// @Failure 404 {object} ErrorResponse
// @Router /pkg/{id} [get]
func getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	var customer model.Customer
	err = db.QueryRow("SELECT id, name, role, email, phone, contacted FROM pkg WHERE id = ?", id).Scan(
		&customer.ID, &customer.Name, &customer.Role, &customer.Email, &customer.Phone, &customer.Contacted)
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
// @Tags pkg
// @Accept json
// @Produce json
// @Param customer body Customer true "Customer to add"
// @Success 201 {object} Customer
// @Router /pkg [post]
func addCustomer(w http.ResponseWriter, r *http.Request) {
	var customer model.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO pkg (name, role, email, phone, contacted) VALUES (?, ?, ?, ?, ?)",
		customer.Name, customer.Role, customer.Email, customer.Phone, customer.Contacted)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	customer.ID = int(id)

	encodeJSONResponse(w, customer)
}

// @Summary Update a customer
// @Description Update a customer
// @Tags pkg
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param customer body Customer true "Customer to update"
// @Success 200 {object} Customer
// @Failure 404 {object} ErrorResponse
// @Router /pkg/{id} [put]
func updateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	var customer model.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE pkg SET name = ?, role = ?, email = ?, phone = ?, contacted = ? WHERE id = ?",
		customer.Name, customer.Role, customer.Email, customer.Phone, customer.Contacted, id)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	customer.ID = id
	encodeJSONResponse(w, customer)
}

// @Summary Delete a customer
// @Description Delete a customer
// @Tags pkg
// @Param id path int true "Customer ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Router /pkg/{id} [delete]
func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM customer WHERE id = ?", id)
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
