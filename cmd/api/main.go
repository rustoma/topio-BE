package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"topio/internal/repository"
	dbrepo "topio/internal/repository/dbRepo"
	"topio/openAI"
)

const port = 8080

type application struct {
	DSN    string
	Domain string
	DB     repository.DatabaseRepo
	AI     ai.AI
}

func main() {
	var app application

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	flag.StringVar(&app.DSN, "dsn", os.Getenv("DSN"), "Postgres connection string")

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	app.DB.CreateTables()
	defer app.DB.Connection().Close()

	app.AI = ai.AI{Client: ai.InitAI()}

	log.Println("Starting application on port", port)

	//start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
