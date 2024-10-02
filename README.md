## FarmApp - Farm Customer Management

### Overview
FarmApp is a simple Go application for managing farm customers. It provides a REST API for performing CRUD operations on customer data and serves static web pages for interacting with the API.

### Features
- Add, view, update, and delete customers.
- API documented with Swagger.
- Static web pages for interacting with the backend.

## Setup Instructions

### 1. Clone the repository
```bash
git clone <repository-url>
cd farmApp
```

### 2. Setting Up Go Environment

1. **Install Go**: Follow the instructions on the [official Go website](https://golang.org/doc/install) to install Go on your machine.
2. **Initialize Go Modules**: Run the following command to initialize Go modules.
    ```bash
    go mod tidy
    ```
3. **Run the Application**: Use the following command to run the application.
    ```bash
    go run main.go
    ```
4. **Generate OpenAPI Documentation**: Use the following command to generate OpenAPI documentation.
    ```bash
    swag init -g main.go --parseDependency
    ```
5. **Serve OpenAPI Documentation**: Ensure the OpenAPI documentation is served by running the application.
    ```bash
    go run main.go
    ```

The application will be available at `http://localhost:8080/swagger/`.

### 3. Docker Setup

1. **Build the Docker image**:
    ```bash
    docker-compose build
    ```

2. **Run the Docker container**:
    ```bash
    docker-compose up
    ```

3. **Stop the Docker container**:
    ```bash
    docker-compose down
    ```

The application will be available at `http://localhost:8080`.

### 4. Accessing the Application

- **Swagger Documentation**: `http://localhost:8080/swagger/`
- **Static Pages**:
    - **Add Customer**: `http://localhost:8080/static/add_customer.html`
    - **View Customers**: `http://localhost:8080/static/view_customers.html`
    - **Update Customer**: `http://localhost:8080/static/update_customer.html`
    - **Delete Customer**: `http://localhost:8080/static/delete_customer.html`

### 5. API Endpoints
- **GET** `/customers` - Retrieve all customers.
- **GET** `/customers/{id}` - Retrieve a customer by ID.
- **POST** `/customers` - Add a new customer.
- **PUT** `/customers/{id}` - Update a customer.
- **DELETE** `/customers/{id}` - Delete a customer.

### 6. Explanation of `index.html`

The `index.html` file provides a user interface for managing farm customers. It includes:
- A form for adding new customers with fields for `Name`, `Role`, `Email`, `Phone`, and `Contacted`.
- A form for updating existing customers with fields for `ID`, `Name`, `Role`, `Email`, `Phone`, and `Contacted`.
- A form for deleting customers by `ID`.
- A table for displaying all customer data with columns for `ID`, `Name`, `Role`, `Email`, `Phone`, and `Contacted`.

### 7. Explanation of OpenAPI

OpenAPI (Swagger) is used to document the API. The documentation is available at `http://localhost:8080/swagger/` and provides a user-friendly interface to interact with the API endpoints. It includes details about the available endpoints, request parameters, and response formats.