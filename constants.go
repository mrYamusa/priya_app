package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var Dbuser = os.Getenv("DbUser")
var Dbpwd = os.Getenv("DbPwd")
var Dbname = os.Getenv("DbName")
