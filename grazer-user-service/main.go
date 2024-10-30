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

func initDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), // should be "mariadb"
		"3306",
		os.Getenv("DB_NAME"))
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the User model
	db.AutoMigrate(&User{})
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if result := db.Create(&user); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var user User
	if result := db.First(&user, id); result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var user User
	if result := db.First(&user, id); result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db.Save(&user)
	json.NewEncoder(w).Encode(user)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	db.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func main() {
	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users", getUsers).Methods("GET")

	fmt.Println("User Service is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
