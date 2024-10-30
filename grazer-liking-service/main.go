// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var rabbitConn *amqp.Connection
var rabbitChannel *amqp.Channel

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

	db.AutoMigrate(&Like{})
}

func initRabbitMQ() {
	var err error
	rabbitConn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	rabbitChannel, err = rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	_, err = rabbitChannel.QueueDeclare(
		"likes", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
}

// Publishes a like event to RabbitMQ
func publishLikeEvent(like Like) error {
	body, err := json.Marshal(like)
	if err != nil {
		return err
	}

	return rabbitChannel.Publish(
		"",      // exchange
		"likes", // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// AddLike creates a like entry for a user and publishes an event to RabbitMQ
func addLike(w http.ResponseWriter, r *http.Request) {
	var like Like
	if err := json.NewDecoder(r.Body).Decode(&like); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	like.CreatedAt = time.Now()
	if result := db.Create(&like); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := publishLikeEvent(like); err != nil {
		log.Printf("Failed to publish like event: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(like)
}

// GetLikes retrieves all likes for a specific user
func getLikes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	var likes []Like
	if result := db.Where("user_id = ?", userID).Find(&likes); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(likes)
}

// GetAllLikes retrieves all likes in the database
func getAllLikes(w http.ResponseWriter, r *http.Request) {
	var likes []Like
	if result := db.Find(&likes); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(likes)
}

func main() {
	initDB()
	initRabbitMQ()
	defer rabbitConn.Close()
	defer rabbitChannel.Close()

	r := mux.NewRouter()
	r.HandleFunc("/likes", addLike).Methods("POST")
	r.HandleFunc("/likes/{userID}", getLikes).Methods("GET")
	r.HandleFunc("/likes", getAllLikes).Methods("GET")

	fmt.Println("Liking Service is running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}
