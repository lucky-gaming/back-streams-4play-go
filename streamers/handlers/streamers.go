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
func GetStreamers(w http.ResponseWriter, r *http.Request) {
	collection := db.DB.Collection("streamers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Erro ao buscar streamers", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var streamers []models.Streamer
	if err = cursor.All(ctx, &streamers); err != nil {
		http.Error(w, "Erro ao processar dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(streamers)
}

func GetStreamerById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	collection := db.DB.Collection("streamers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.Streamer
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		http.Error(w, "Moderador não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// POST /streamers
func CreateStreamer(w http.ResponseWriter, r *http.Request) {
	var streamer models.Streamer
	if err := json.NewDecoder(r.Body).Decode(&streamer); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	streamer.ID = primitive.NewObjectID()

	collection := db.DB.Collection("streamers")
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

func UpdateStreamer(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	objectId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Remove campos que não devem ser atualizados
	delete(payload, "id")
	delete(payload, "_id")

	if len(payload) == 0 {
		http.Error(w, "Nenhum campo para atualizar", http.StatusBadRequest)
		return
	}

	collection := db.DB.Collection("streamers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": payload}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectId}, update)
	if err != nil {
		http.Error(w, "Erro ao atualizar Streamer", http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		http.Error(w, "Streamer não encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"message": "Atualizado com sucesso"})
}


// handlers/password.go
func DeleteStreamer(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	objectId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	collection := db.DB.Collection("streamers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil || result.DeletedCount == 0 {
		http.Error(w, "Erro ao deletar streamer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"message": "Deletado com sucesso"})
}
