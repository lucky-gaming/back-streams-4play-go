package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go-api/db"
	"go-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GET /users
func GetModerators(w http.ResponseWriter, r *http.Request) {
	collection := db.DB.Collection("moderators")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Erro ao buscar moderators", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var moderators []models.Moderator
	if err = cursor.All(ctx, &moderators); err != nil {
		http.Error(w, "Erro ao processar dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(moderators)
}

func GetModeratorById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	collection := db.DB.Collection("moderators")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.Moderator
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		http.Error(w, "Moderador não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// POST /moderators
func CreateModerator(w http.ResponseWriter, r *http.Request) {
	var streamer models.Moderator
	if err := json.NewDecoder(r.Body).Decode(&streamer); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	streamer.ID = primitive.NewObjectID()

	collection := db.DB.Collection("moderators")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, streamer)
	if err != nil {
		http.Error(w, "Erro ao salvar streamer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(streamer)
}

func UpdateModerator(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	objectId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Decodifica o JSON como um map genérico
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Remove campos proibidos de serem atualizados diretamente
	delete(payload, "id")
	delete(payload, "_id")

	// Verifica se há algo a ser atualizado
	if len(payload) == 0 {
		http.Error(w, "Nenhum campo para atualizar", http.StatusBadRequest)
		return
	}

	collection := db.DB.Collection("moderators")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": payload}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectId}, update)
	if err != nil {
		http.Error(w, "Erro ao atualizar Moderator", http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		http.Error(w, "Moderator não encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"message": "Atualizado com sucesso"})
}


// handlers/password.go
func DeleteModerator(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	objectId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	collection := db.DB.Collection("moderators")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil || result.DeletedCount == 0 {
		http.Error(w, "Erro ao deletar moderador", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"message": "Deletado com sucesso"})
}
