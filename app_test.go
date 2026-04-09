package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var app App

func TestMain(m *testing.M) {
	app.Initialise(os.Getenv("Dbuser"), os.Getenv("Dbpwd"), "test")
	createTable()

	code := m.Run() // run tests

	os.Exit(code) // VERY IMPORTANT
}

func createTable() {
	clearTable()

	createTableQuery := `CREATE TABLE IF NOT EXISTS products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		quantity INT NOT NULL,
		price DECIMAL(10, 2) NOT NULL
	)`

	_, err := app.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Error while creating products table -", err)
	}
}

func clearTable() {
	_, err := app.DB.Exec("DELETE FROM products")

	if err != nil {
		log.Println("Error while clearing products table - ", err)
		log.Fatal("Error while clearing products table - ", err)
	}
	app.DB.Exec("ALTER TABLE products AUTO_INCREMENT = 1")
}

func addProducts(name string, quantity int, price float64) {
	insertQuery := `INSERT INTO products (name, quantity, price) VALUES (?, ?, ?)`
	_, err := app.DB.Exec(insertQuery, name, quantity, price)
	if err != nil {
		log.Fatal("Error while adding product -", err)
	}
	log.Printf("Product added successfully - Name: %v, Quantity: %v, Price: %v\n", name, quantity, price)
}

func TestGetProduct(t *testing.T) {
	log.Println("Testing GetProduct endpoint")
	clearTable()
	log.Println("Cleared products table")
	log.Println("Adding sample products to the database")
	addProducts("Product 1", 10, 99.99)
	addProducts("Product 2", 5, 49.99)
	log.Println("Sample products added successfully")
	req, err := http.NewRequest("GET", "/product/1", nil)
	if err != nil {
		log.Panicln("Error while creating request -", err)
		t.Fatal("Error while creating request -", err)
	}
	recorder := sendRequest(req)
	checkStatusCode(t, http.StatusAccepted, recorder)
	log.Println("GetProduct endpoint test passed successfully")
}

func sendRequest(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	app.Router.ServeHTTP(recorder, req)
	log.Printf("Request sent to %v with method %v\n", req.URL.Path, req.Method)
	log.Printf("Response received with status code %v\n", recorder.Code)
	return recorder
}

func checkStatusCode(t *testing.T, expected int, recorder *httptest.ResponseRecorder) {
	log.Printf("Checking if expected status code %v matches actual status code %v\n", expected, recorder.Code)
	if recorder.Code != expected {
		log.Printf("Status code check failed - expected %v but got %v\n", expected, recorder.Code)
		t.Errorf("Expected status code %v but got %v\n", expected, recorder.Code)
	}
	log.Printf("Status code check passed - expected %v matches actual %v\n", expected, recorder.Code)
}

// http serve
// http new request
// httptest new recorder
