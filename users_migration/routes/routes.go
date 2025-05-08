package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go-api/handlers"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // ou especifique: http://localhost:3000
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutos
	}))

	r.Post("/gtm_info", handlers.CreateGtmInfo)

	return r
}
