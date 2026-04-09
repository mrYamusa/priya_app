package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	app := &App{}
	app.Initialise(Dbuser, Dbpwd, Dbname)
	log.Println("App starting on port localhost:8080")
	app.Run()
}
