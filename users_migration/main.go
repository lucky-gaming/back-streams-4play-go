package main

import (
	"log"
	"net/http"
	"go-api/db"
	"go-api/routes"
)

func main() {
	db.Connect()

	r := routes.SetupRoutes()

	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
