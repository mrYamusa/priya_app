package main

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	app := &App{}
	app.Initialise(os.Getenv("Dbuser"), os.Getenv("Dbpwd"), os.Getenv("Dbname"))
	log.Println("App starting on port localhost:8080")
	app.Run()
}
