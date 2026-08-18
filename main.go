package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	_ "github.com/lib/pq"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Hash struct {
	Salt           string
	HashedPassword string
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	http.HandleFunc("/users", createUserHandler)

	fmt.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func generateSalt() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "0123456789abcdef"
	b := make([]byte, 12)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateHash(user User, salt string) Hash {
	mac := hmac.New(sha512.New, []byte(salt))
	mac.Write([]byte(user.Password))
	return Hash{
		Salt:           salt,
		HashedPassword: hex.EncodeToString(mac.Sum(nil)),
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Printf("error connecting to db: %v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hash := generateHash(user, generateSalt())
	id := fmt.Sprintf("%x", rand.Int63())

	_, err = db.Exec(`
		INSERT INTO users (id, email, password_hash, salt, checked, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, user.Email, hash.HashedPassword, hash.Salt, false, time.Now(), time.Now(),
	)
	if err != nil {
		log.Printf("error inserting user: %v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if err := publishEvent(user); err != nil {
		log.Printf("failed to publish user event: %v", err)
	}

	log.Printf("Successfully created user: %+v", user)

	fmt.Fprintf(w, `{"id":"%s"}`, id)
	w.WriteHeader(http.StatusCreated)
}

func publishEvent(user User) error {
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, _ := json.Marshal(user)
	return ch.Publish(
		"",
		"users.new",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
