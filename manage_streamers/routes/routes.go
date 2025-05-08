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

	r.Get("/streamers", handlers.GetStreamers)
	r.Get("/streamers/{id}", handlers.GetStreamerById)
	r.Post("/streamers", handlers.CreateStreamer)
	r.Put("/streamers/{id}", handlers.UpdateStreamer)
	r.Delete("/streamers/{id}", handlers.DeleteStreamer)


	r.Get("/moderators", handlers.GetModerators)
	r.Get("/moderators/{id}", handlers.GetModeratorById)
	r.Post("/moderators", handlers.CreateModerator)
	r.Put("/moderators/{id}", handlers.UpdateModerator)
	r.Delete("/moderators/{id}", handlers.DeleteModerator)

	r.Get("/lives", handlers.GetLivesByStatus)
	r.Get("/lives/{id}", handlers.GetLiveById)
	r.Post("/lives", handlers.CreateLive)
	r.Put("/lives/{id}", handlers.UpdateLive)
	r.Delete("/lives/{id}", handlers.DeleteLive)
	r.Get("/lives/active_at", handlers.CheckLiveAtTimestamp)

	return r
}
