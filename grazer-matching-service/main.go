// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

	db.AutoMigrate(&Match{})
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
}

func processLikes() {
	msgs, err := rabbitChannel.Consume(
		"likes", // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for d := range msgs {
		var like Like
		if err := json.Unmarshal(d.Body, &like); err != nil {
			log.Printf("Error decoding message: %v", err)
			continue
		}
		checkAndCreateMatch(like)
	}
}

func checkAndCreateMatch(like Like) {
	var reciprocalLike Like
	reciprocalExists := db.Raw("SELECT * FROM likes WHERE user_id = ? AND liked_user_id = ?", like.LikedUserID, like.UserID).
		Scan(&reciprocalLike).RowsAffected > 0

	if reciprocalExists {
		match := Match{
			UserID:        like.UserID,
			MatchedUserID: like.LikedUserID,
			CreatedAt:     time.Now(),
		}
		db.Create(&match)
		log.Printf("Match created between %d and %d", like.UserID, like.LikedUserID)
	}
}

// AddMatch creates a new match manually (for testing)
func addMatch(w http.ResponseWriter, r *http.Request) {
	var match Match
	if err := json.NewDecoder(r.Body).Decode(&match); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match.CreatedAt = time.Now()
	if result := db.Create(&match); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(match)
}

// GetMatches retrieves all matches for a specific user
func getMatches(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["userID"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var matches []Match
	db.Where("user_id = ? OR matched_user_id = ?", userID, userID).Find(&matches)

	json.NewEncoder(w).Encode(matches)
}

// GetAllMatches retrieves all matches in the database
func getAllMatches(w http.ResponseWriter, r *http.Request) {
	var matches []Match
	if result := db.Find(&matches); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(matches)
}

func main() {
	initDB()
	initRabbitMQ()
	defer rabbitConn.Close()
	defer rabbitChannel.Close()

	go processLikes()

	r := mux.NewRouter()
	r.HandleFunc("/matches", addMatch).Methods("POST") // New POST endpoint for manual match creation
	r.HandleFunc("/matches/{userID}", getMatches).Methods("GET")
	r.HandleFunc("/matches", getAllMatches).Methods("GET")

	fmt.Println("Matching Service is running and listening for likes...")
	log.Fatal(http.ListenAndServe(":8083", r))
}
