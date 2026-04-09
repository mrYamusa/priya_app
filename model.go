package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func GetProducts(s *sql.DB) ([]Product, error) {
	rows, err := s.Query("SELECT * from products")
	if err != nil {
		return nil, err
	}

	products := []Product{}
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, err
}

func (p *Product) GetProduct(s *sql.DB) error {
	err := s.QueryRow("SELECT * from products where id=?", p.ID).Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			log.Printf("Product with ID: %v not found!\n", p.ID)
			return fmt.Errorf("Product with ID: %v not found!", p.ID)
		default:
			log.Printf("Error while fetching product with ID: %v - %v\n", p.ID, err)
			return err
		}
	}
	return nil
}

func (p *Product) MakeProduct(s *sql.DB, r *http.Request) (Product, error) {
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		return Product{}, err
	}
	log.Printf("Creating product with name: %v, quantity: %v and price: %v\n", p.Name, p.Quantity, p.Price)
	result, err := s.Exec("INSERT INTO products(name, quantity, price) VALUES(?, ?, ?)", p.Name, p.Quantity, p.Price)
	if err != nil {
		log.Printf("Error while creating product - %v\n", err)
		return Product{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error while fetching last inserted id - %v\n", err)
		return Product{}, err
	}
	p.ID = int(id)
	return *p, nil
}

func (p *Product) UpdateProducts(s *sql.DB) error {

	log.Printf("Updating product with ID: %v\n", p.ID)
	_, err := s.Exec("update products set name=?, quantity=?, price=? where id=?", p.Name, p.Quantity, p.Price, p.ID)
	if err != nil {
		log.Printf("Error while updating product - %v\n", err)
		return err
	}
	log.Printf("Product with ID: %v updated successfully\n", p.ID)
	return nil
}

func (p *Product) DeleteProduct(s *sql.DB) error {
	result, err := s.Exec("delete from products where id=?", p.ID)
	if err != nil {
		log.Println("Error deleting product with id: ", p.ID)
		return err
	}
	var n int64
	if n, err = result.RowsAffected(); n == 0 {
		log.Printf("Product with id: %v doesn't exist", p.ID)
		return fmt.Errorf("Product with id: %v doesn't exist", p.ID)
	}
	return nil
}
