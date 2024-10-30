// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// Initialize the database connection
func initDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		"3306",
		os.Getenv("DB_NAME"))

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the Preferences model
	db.AutoMigrate(&Preferences{})
}

// Handlers for the API endpoints

// createPreferences creates new preferences for a user
func createPreferences(w http.ResponseWriter, r *http.Request) {
	var preferences Preferences
	if err := json.NewDecoder(r.Body).Decode(&preferences); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if result := db.Create(&preferences); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(preferences)
}

// getPreferences retrieves preferences for a specific user
func getPreferences(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	var preferences Preferences
	if result := db.First(&preferences, "user_id = ?", userID); result.Error != nil {
		http.Error(w, "Preferences not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(preferences)
}

// getAllPreferences retrieves all user preferences
func getAllPreferences(w http.ResponseWriter, r *http.Request) {
	var preferences []Preferences
	if result := db.Find(&preferences); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(preferences)
}

// updatePreferences updates existing user preferences
func updatePreferences(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	var preferences Preferences
	if result := db.First(&preferences, "user_id = ?", userID); result.Error != nil {
		http.Error(w, "Preferences not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&preferences); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db.Save(&preferences)
	json.NewEncoder(w).Encode(preferences)
}

func main() {
	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/preferences", createPreferences).Methods("POST")
	r.HandleFunc("/preferences/{userID}", getPreferences).Methods("GET")
	r.HandleFunc("/preferences", getAllPreferences).Methods("GET") // New endpoint
	r.HandleFunc("/preferences/{userID}", updatePreferences).Methods("PUT")

	fmt.Println("Preferences Service is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
