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

	var Dbuser = os.Getenv("Dbuser")
	var Dbpwd = os.Getenv("Dbpwd")
	var Dbname = os.Getenv("Dbname")

	log.Println("Database credentials loaded successfully - User:", Dbuser, "Password:", Dbpwd, "Database Name:", Dbname)
}
