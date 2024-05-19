package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux" // Importing mux package

	"go.mongodb.org/mongo-driver/bson"

	"backend_golang/db"
	"backend_golang/models"
	"backend_golang/logger"
)

func InsertDataHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var message models.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Log("Error decoding request body:", err)
		return
	}

	// Generate unique random 8-digit alphanumeric ID
	for {
		message.ID = generateRandomID(8)
		filter := bson.M{"id": message.ID}
		count, err := db.MessagesCollection.CountDocuments(context.Background(), filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Log("Error checking for existing ID in MongoDB:", err)
			return
		}
		if count == 0 {
			break
		}
		// ID already exists, generate a new one
	}

	// Insert data into MongoDB
	result, err := db.MessagesCollection.InsertOne(context.Background(), message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Log("Error inserting data into MongoDB:", err)
		return
	}

	// Get the inserted document
	var insertedMessage models.Message
	err = db.MessagesCollection.FindOne(context.Background(), bson.M{"_id": result.InsertedID}).Decode(&insertedMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Log("Error retrieving inserted document:", err)
		return
	}

	// Log inserted data with timestamp
	logger.Printf("[%s] Data inserted: %+v\n", time.Now().Format(time.Stamp), insertedMessage)

	// Respond with the inserted document
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Data inserted successfully")
	logger.Log("Data inserted successfully")
	json.NewEncoder(w).Encode(insertedMessage)
}

func GetDataHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve data from MongoDB
	cursor, err := db.MessagesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Log("Error retrieving data from MongoDB:", err)
		return
	}
	defer cursor.Close(context.Background())

	// Extract data from cursor
	var messages []models.Message
	for cursor.Next(context.Background()) {
		var message models.Message
		err := cursor.Decode(&message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Log("Error decoding message from cursor:", err)
			return
		}
		messages = append(messages, message)
	}

	// Convert result to JSON
	jsonData, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Log("Error marshalling messages to JSON:", err)
		return
	}

	// Set Content-Type header and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	logger.Log("Data retrieved successfully")
}

func UpdateDataHandler(w http.ResponseWriter, r *http.Request) {
	// Parse ID parameter
	params := mux.Vars(r)
	id := params["id"]

	// Parse request body
	var updatedMessage models.Message
	if err := json.NewDecoder(r.Body).Decode(&updatedMessage); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Log("Error decoding request body:", err)
		return
	}

	// Update data in MongoDB
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{
		"name":    updatedMessage.Name,
		"email":   updatedMessage.Email,
		"message": updatedMessage.Message,
		"status":  updatedMessage.Status,
	}}
	result, err := db.MessagesCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Log("Error updating data in MongoDB:", err)
		return
	}

	// Check if any document was updated
	if result.ModifiedCount == 0 {
		http.Error(w, fmt.Sprintf("No document found with ID: %s", id), http.StatusNotFound)
		logger.Printf("No document found with ID: %s\n", id)
		return
	}
	// Get the updated document
	var updatedDoc models.Message
	err = db.MessagesCollection.FindOne(context.Background(), bson.M{"id": id}).Decode(&updatedDoc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Log("Error retrieving updated document:", err)
		return
	}

	// Log updated data with timestamp
	logger.Printf("[%s] Data updated: %+v\n", time.Now().Format(time.Stamp), updatedMessage)

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	logger.Log("Data updated successfully")
	json.NewEncoder(w).Encode(updatedDoc)
	fmt.Fprintf(w, "Data updated successfully")
}

func DeleteDataHandler(w http.ResponseWriter, r *http.Request) {
	// Parse ID parameter
	params := mux.Vars(r)
	id := params["id"]

	// Delete data from MongoDB
	filter := bson.M{"id": id}
	_, err := db.MessagesCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Log("Error deleting data from MongoDB:", err)
		return
	}

	// Log deletion with timestamp
	logger.Printf("[%s] Data deleted with ID: %s\n", time.Now().Format(time.Stamp), id)

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Data deleted successfully")
	logger.Log("Data deleted successfully")
}

func PatchDataHandler(w http.ResponseWriter, r *http.Request) {
	// Parse ID parameter
	params := mux.Vars(r)
	id := params["id"]

	// Parse request body
	var statusUpdate struct {
		Status bool `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Log("Error decoding request body:", err)
		return
	}

	// Update data in MongoDB
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"status": statusUpdate.Status}}
	result, err := db.MessagesCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Log("Error updating data in MongoDB:", err)
		return
	}

	// Check if any document was updated
	if result.ModifiedCount == 0 {
		http.Error(w, fmt.Sprintf("No document found with ID: %s", id), http.StatusNotFound)
		logger.Printf("No document found with ID: %s\n", id)
		return
	}

	// Get the updated document
	var updatedDoc models.Message
	err = db.MessagesCollection.FindOne(context.Background(), bson.M{"id": id}).Decode(&updatedDoc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Log("Error retrieving updated document:", err)
		return
	}

	// Log updated data with timestamp
	logger.Printf("[%s] Data status updated: %+v\n", time.Now().Format(time.Stamp), updatedDoc)

	// Respond with success message
	w.WriteHeader(http.StatusOK)
	logger.Log("Data status updated successfully")
	json.NewEncoder(w).Encode(updatedDoc)
	fmt.Fprintf(w, "Data status updated successfully")
}	