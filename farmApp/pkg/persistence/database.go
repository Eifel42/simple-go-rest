package persistence

import (
	"database/sql"
	"farmApp/pkg/api"
	"log"
	"os"
)

const databaseName = "./farmCustomers.db"

var db *sql.DB

func CreateFarmDB() {
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
    CREATE TABLE IF NOT EXISTS customer (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        role TEXT,
        email TEXT,
        phone TEXT,
        contacted BOOLEAN
    );`
	return execQuery(query)
}

func execQuery(query string) error {
	_, err := db.Exec(query)
	return err
}

func insertInitialCustomers() error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM customer").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		log.Println("No customers found. Inserting initial data...")
		customers := []api.Customer{
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

func bulkInsertCustomers(customers []api.Customer) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO customer (name, role, email, phone, contacted) VALUES (?, ?, ?, ?, ?)")
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

func GetCustomers() ([]api.Customer, error) {
	rows, err := db.Query("SELECT id, name, role, email, phone, contacted FROM customer")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []api.Customer
	for rows.Next() {
		var customer api.Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Role, &customer.Email, &customer.Phone, &customer.Contacted); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func GetCustomerByID(id int) (api.Customer, error) {
	var customer api.Customer
	err := db.QueryRow("SELECT id, name, role, email, phone, contacted FROM customer WHERE id = ?", id).Scan(
		&customer.ID, &customer.Name, &customer.Role, &customer.Email, &customer.Phone, &customer.Contacted)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func AddCustomer(customer api.Customer) (int, error) {
	result, err := db.Exec("INSERT INTO customer (name, role, email, phone, contacted) VALUES (?, ?, ?, ?, ?)",
		customer.Name, customer.Role, customer.Email, customer.Phone, customer.Contacted)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func UpdateCustomer(id int, customer api.Customer) error {
	_, err := db.Exec("UPDATE customer SET name = ?, role = ?, email = ?, phone = ?, contacted = ? WHERE id = ?",
		customer.Name, customer.Role, customer.Email, customer.Phone, customer.Contacted, id)
	return err
}

func DeleteCustomer(id int) error {
	_, err := db.Exec("DELETE FROM customer WHERE id = ?", id)
	return err
}
