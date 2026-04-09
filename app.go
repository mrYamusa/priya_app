package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func checkErr(err error) {
	if err != nil {
		log.Println("The error encountered:", err)
	}
}

func (app *App) Initialise(Dbuser string, Dbpwd string, Dbname string) {
	log.Println("Initialising the app")
	// Router
	app.Router = mux.NewRouter().StrictSlash(true)

	// DataBase
	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", Dbuser, Dbpwd, Dbname)
	var err error
	app.DB, err = sql.Open("mysql", connectionString)
	os.Getenv("DbUser")
	checkErr(err)

	// Make Routes
	app.HandleRoutes()
	log.Println("App initialised successfully")
}

func (app *App) Run() {
	log.Fatal(http.ListenAndServe(":8080", app.Router))
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	log.Printf("Sending response with status code %v\n", statusCode)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	response, _ := json.Marshal(payload)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	log.Printf("Sending error response with status code %v and error message %v\n", statusCode, err)
	error_message := map[string]string{"error": err}
	sendResponse(w, statusCode, error_message)
}

// func getProducts(s *sql.DB) {

// }

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting all products form /products")
	products, err := GetProducts(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusAccepted, products)
	log.Println("Products retrieved successfully")
}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating product form /product")
	var p Product
	item, err := p.MakeProduct(app.DB, r)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusCreated, item)
	log.Println("Product created successfully")
}

// insert := fmt.Sprintf("insert into products values(%v, %v, %v)", )

func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting product form /product/{id}")
	v := mux.Vars(r)
	a, err := strconv.Atoi(v["id"])
	if err != nil {
		log.Println("Could not parse path parameter - ", err)
		sendError(w, http.StatusNotAcceptable, err.Error())
		return
	}
	var p Product = Product{ID: a}
	err = p.GetProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusNotFound, err.Error())
		return
	}
	sendResponse(w, http.StatusAccepted, p)
	log.Printf("Product with ID: %v retrieved successfully\n", p.ID)

}
func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating product form /product/{id}")
	v := mux.Vars(r)
	a, err := strconv.Atoi(v["id"])
	if err != nil {
		log.Printf("Could not parse parameter var: %v\n", v["id"])
	}

	var p Product = Product{ID: a}
	err = p.GetProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusNotFound, fmt.Sprintf("Product with ID: %v not found!", a))
		return
	}
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Printf("Error while decoding product data: %v\n", err)
		sendError(w, http.StatusBadRequest, "Invalid product data")
		return
	}
	err = p.UpdateProducts(app.DB)
	if err != nil {
		log.Printf("Error while updating product with ID: %v\n", p.ID)
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusAccepted, p)
	log.Printf("Product with ID: %v updated successfully\n", p.ID)
}
func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	a := mux.Vars(r)
	b, err := strconv.Atoi(a["id"])
	if err != nil {
		log.Println("Invalid Path parameter: ", a["id"])
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	var p = Product{ID: b}
	err = p.DeleteProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusNotFound, err.Error())
		return
	}
	sendResponse(w, http.StatusAccepted, p.ID)
	log.Printf("Product with ID: %v deleted successfully\n", p.ID)
}

func (app *App) HandleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")

	app.Router.HandleFunc("/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")

	app.Router.HandleFunc("/reference", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "https://generator3.swagger.io/openapi.json", // allow external URL or local path file
			// SpecURL: "./docs/swagger.json",
			Theme: scalar.ThemeDeepSpace,
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Simple API",
			},
			DarkMode: true,
		})

		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Fprintln(w, htmlContent)
	}).Methods("GET")
}
