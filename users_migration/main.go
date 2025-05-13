package main

import (
	"go-api/db"
	"go-api/routes"
	"log"
	"net/http"
	"os"
)

func main() {

	db.Connect()

	r := routes.SetupRoutes()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("APP_HOST")
	//Habilitar port segura
	if os.Getenv("HTTPS_ENABLED") == "true" {
		log.Println("Iniciando o servidor em ", host, " porta", port, "com TLS...")
		// Load the TLS certificate and key
		certFile := "cert.pem"
		keyFile := "key.pem"

		err := http.ListenAndServeTLS(":"+port, certFile, keyFile, r)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		//Habilitar port não segura
		err := http.ListenAndServe(":"+port, r)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Iniciando o servidor em ", host, " porta", port, "com TLS...")
	}
}
