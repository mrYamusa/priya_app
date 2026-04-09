package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func GetProductId(basepath string, s string) (i int, err error) {
	a := len(basepath)
	r := s[a:]
	z, err := strconv.Atoi(r)
	if err != nil {
		log.Println("couldn't coonvert to string: ", err)
	}

	fmt.Printf("%v\n%+v\n%T\n", z, z, z)
	return z, err
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GetProductId("/products/", "/products/20")
	fmt.Println(os.Getenv("DbUser"))
}
