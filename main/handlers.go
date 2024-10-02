package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)

const databaseName = "./customerDatabase.db"

var db *sql.DB

type Customer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Contacted bool   `json:"contacted"`
}

func createDB() {
	deleteDB(databaseName)
	var err error
	db, err = sql.Open("sqlite3", databaseName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	err = initDB(databaseName)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteDB(dataSourceName string) {
	err := os.Remove(dataSourceName)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Failed to delete the database:", err)
	}
	log.Println("Database deleted successfully.")
}

func initDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}

	if err = createCustomersTable(); err != nil {
		return err
	}

	if err = insertInitialCustomers(); err != nil {
		return err
	}

	return nil
}

func createCustomersTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS customers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        role TEXT,
        email TEXT,
        phone TEXT,
        contacted BOOLEAN
    );`
	return execQuery(query)
}

func insertInitialCustomers() error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM customers").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		log.Println("No customers found. Inserting initial data...")
		customers := []Customer{
			{Name: "Bauer Klaus", Role: "Farmer", Email: "klaus.bauer@farm.de", Phone: "01234 567890", Contacted: true},
			{Name: "Bauerin Anna", Role: "Owner", Email: "anna.bauerin@farm.de", Phone: "01234 567891", Contacted: false},
			{Name: "Müller Hans", Role: "Worker", Email: "hans.mueller@farm.de", Phone: "01234 567892", Contacted: true},
			{Name: "Schmidt Peter", Role: "Manager", Email: "peter.schmidt@farm.de", Phone: "01234 567893", Contacted: false},
			{Name: "Fischer Maria", Role: "Assistant", Email: "maria.fischer@farm.de", Phone: "01234 567894", Contacted: true},
			{Name: "Weber Karl", Role: "Technician", Email: "karl.weber@farm.de", Phone: "01234 567895", Contacted: false},
			{Name: "Meyer Lisa", Role: "Accountant", Email: "lisa.meyer@farm.de", Phone: "01234 567896", Contacted: true},
			{Name: "Wagner Thomas", Role: "Driver", Email: "thomas.wagner@farm.de", Phone: "01234 567897", Contacted: false},
			{Name: "Becker Laura", Role: "Secretary", Email: "laura.becker@farm.de", Phone: "01234 567898", Contacted: true},
			{Name: "Hoffmann Frank", Role: "Guard", Email: "frank.hoffmann@farm.de", Phone: "01234 567899", Contacted: false},
		}

		err := bulkInsertCustomers(customers)
		if err != nil {
			return err
		}
		log.Println("Inserted initial customers.")
	} else {
		log.Println("Customers already exist in the database.")
	}
	return nil
}

func insertCustomer(customer Customer) error {
	query := "INSERT INTO customers (name, role, email, phone, contacted) VALUES (?, ?, ?, ?, ?)"
	return execStmtWithTransaction(query, customer.Name, customer.Role, customer.Email, customer.Phone, customer.Contacted)
}

func updateCustomerDB(id string, customer Customer) error {
	query := "UPDATE customers SET name = ?, role = ?, email = ?, phone = ?, contacted = ? WHERE id = ?"
	return execStmtWithTransaction(query, customer.Name, customer.Role, customer.Email, customer.Phone, customer.Contacted, id)
}

func deleteCustomerDB(id string) error {
	query := "DELETE FROM customers WHERE id = ?"
	return execStmtWithTransaction(query, id)
}

func bulkInsertCustomers(customers []Customer) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO customers (name, role, email, phone, contacted) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	for _, customer := range customers {
		_, err = stmt.Exec(customer.Name, customer.Role, customer.Email, customer.Phone, customer.Contacted)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return err
		}
	}

	return tx.Commit()
}

func execQuery(query string) error {
	_, err := db.Exec(query)
	return err
}

func execStmtWithTransaction(query string, args ...interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	_, err = stmt.Exec(args...)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func handleError(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
}

func encodeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Failed to encode JSON response:", err)
	}
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, role, email, phone, contacted FROM customers")
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var customers []Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Role, &customer.Email, &customer.Phone, &customer.Contacted); err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		customers = append(customers, customer)
	}

	if handleRowsError(w, rows.Err()) {
		return
	}

	encodeJSONResponse(w, customers)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var customer Customer
	err := db.QueryRow("SELECT id, name, role, email, phone, contacted FROM customers WHERE id = ?", id).Scan(
		&customer.ID, &customer.Name, &customer.Role, &customer.Email, &customer.Phone, &customer.Contacted)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Customer not found", http.StatusNotFound)
		} else {
			handleError(w, err, http.StatusInternalServerError)
		}
		return
	}

	encodeJSONResponse(w, customer)
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	if err := insertCustomer(customer); err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encodeJSONResponse(w, customer)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	if err := updateCustomerDB(id, customer); err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	encodeJSONResponse(w, customer)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := deleteCustomerDB(id); err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func OpenAPIHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("openapi.yaml")
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	if _, err := w.Write(data); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func handleRowsError(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}
