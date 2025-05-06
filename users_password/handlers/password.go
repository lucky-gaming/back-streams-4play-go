package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go-api/db"
	"go-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// POST /users
func CreateGtmInfo(w http.ResponseWriter, r *http.Request) {
	var user models.UserInfo
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	user.ID = primitive.NewObjectID()

	collection := db.DB.Collection("users_password")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Erro ao salvar usuário", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
