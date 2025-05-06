package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"strconv"

	"go-api/db"
	"go-api/models"
	"go-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GET /lives
func GetLives(w http.ResponseWriter, r *http.Request) {
	collection := db.DB.Collection("lives")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Erro ao buscar lives", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var lives []models.Live
	if err = cursor.All(ctx, &lives); err != nil {
		http.Error(w, "Erro ao processar dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lives)
}

func GetLiveById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	collection := db.DB.Collection("lives")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var live models.Live
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&live)
	if err != nil {
		http.Error(w, "Moderador não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(live)
}

// POST /lives
func CreateLive(w http.ResponseWriter, r *http.Request) {
	var live models.Live
	if err := json.NewDecoder(r.Body).Decode(&live); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	live.ID = primitive.NewObjectID()
	
	live.StreamKey = utils.GenerateStreamKey()

	collection := db.DB.Collection("lives")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, live)
	if err != nil {
		http.Error(w, "Erro ao salvar live", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(live)
}

func UpdateLive(w http.ResponseWriter, r *http.Request) {
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

	collection := db.DB.Collection("lives")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": payload}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectId}, update)
	if err != nil {
		http.Error(w, "Erro ao atualizar Live", http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		http.Error(w, "Live não encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"message": "Atualizado com sucesso"})
}


// handlers/password.go
func DeleteLive(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	objectId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	collection := db.DB.Collection("lives")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil || result.DeletedCount == 0 {
		http.Error(w, "Erro ao deletar live", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"message": "Deletado com sucesso"})
}

func GetLivesByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	start_at := r.URL.Query().Get("start_at")
	

	collection := db.DB.Collection("lives")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if status != "" {
		filter["status"] = status

		if status == "scheduled" {
			var timestamp int64 ;

			if start_at != "" {
				val,err := strconv.ParseInt(start_at, 10, 64)
				if err != nil {
					http.Error(w, "Timestamp inválido", http.StatusBadRequest)
					return
				}
				timestamp = val
			} else {
				timestamp = time.Now().Unix()
			}
			
			filter["startDate"] = bson.M{"$gte": timestamp}
		}


	}

	

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Erro ao buscar lives", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var lives []models.Live
	if err := cursor.All(ctx, &lives); err != nil {
		http.Error(w, "Erro ao decodificar lives", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lives)
}


func CheckLiveAtTimestamp(w http.ResponseWriter, r *http.Request) {
	startAtStr := r.URL.Query().Get("start_at")
	endAtStr := r.URL.Query().Get("end_at")

	if startAtStr == "" || endAtStr == "" {
		http.Error(w, "Parâmetros 'start_at' e 'end_at' são obrigatórios", http.StatusBadRequest)
		return
	}

	startAt, errStart := strconv.ParseInt(startAtStr, 10, 64)
	endAt, errEnd := strconv.ParseInt(endAtStr, 10, 64)

	if errStart != nil || errEnd != nil {
		http.Error(w, "Timestamps inválidos", http.StatusBadRequest)
		return
	}

	collection := db.DB.Collection("lives")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verifica interseção: start_at < live.endDate AND end_at > live.startDate
	filter := bson.M{
		"$and": []bson.M{
			{"endDate": bson.M{"$gt": startAt}},
			{"startDate": bson.M{"$lt": endAt}},
		},
	}

	var live models.Live
	err := collection.FindOne(ctx, filter).Decode(&live)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bson.M{"conflict": false})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bson.M{
		"conflict": true,
		"live":     live,
	})
}



