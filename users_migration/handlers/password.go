package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go-api/db"
	"go-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// POST /users
func CreateGtmInfo(w http.ResponseWriter, r *http.Request) {
	var user models.UserInfo
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if user.SmarticoUserId == "" {
		http.Error(w, "Campo 'smartico_user_id' é obrigatório", http.StatusBadRequest)
		return
	}

	collection := db.DB.Collection("users_password")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verifica se já existe um usuário com o mesmo smartico_user_id
	var existing models.UserInfo
	err := collection.FindOne(ctx, bson.M{"smarticouserid": user.SmarticoUserId}).Decode(&existing)
	if err == nil {
		http.Error(w, "Tracked", http.StatusConflict)
		return
	}

	user.ID = primitive.NewObjectID()

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Error on Track", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bson.M{"message": "Success"})
}
