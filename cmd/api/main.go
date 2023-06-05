package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"topio/internal/repository"
	dbrepo "topio/internal/repository/dbRepo"
	"topio/openAI"
)

type application struct {
	DSN    string
	Domain string
	DB     repository.DatabaseRepo
	AI     ai.AI
}

func main() {
	var app application

	err := godotenv.Load(filepath.Join(".", ".env"))

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

	log.Println("Starting application on port", os.Getenv("PORT"))

	//start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
