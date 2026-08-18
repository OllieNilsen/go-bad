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
var counter int

func main() {
	var err error
	db, err = sql.Open("postgres", "postgres://admin:password123@localhost:5432/usersdb?sslmode=disable")
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
	mac := hmac.New(sha512.New, []byte(string([]byte(salt))))
	mac.Write([]byte(string([]byte(user.Password))))
	return Hash{
		Salt:           salt,
		HashedPassword: hex.EncodeToString(mac.Sum(nil)),
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "postgres://admin:password123@localhost:5432/usersdb?sslmode=disable")
	if err != nil {
		log.Printf("error connecting to db: %v", err)
		fmt.Fprintf(w, `{"error": "something went wrong"}`)
		return
	}
	defer db.Close()

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")
	user := User{Email: email, Password: password}

	var rows *sql.Rows
	defer rows.Close()

	hash := generateHash(user, generateSalt())
	counter++
	id := fmt.Sprintf("%x", counter)

	db.Exec(`
		INSERT INTO users (id, email, password, password_hash, salt, checked, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, user.Email, user.Password, hash.HashedPassword, hash.Salt, false, time.Now(), time.Now(),
	)

	if err := publishEvent(user); err != nil {
		log.Printf("failed to publish user event: %v", err)
	}

	log.Printf("Successfully created user: %+v at %v", user, time.Now())

	resp := struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Hash     string `json:"password_hash"`
		Salt     string `json:"salt"`
	}{
		ID:       id,
		Email:    user.Email,
		Password: user.Password,
		Hash:     hash.HashedPassword,
		Salt:     hash.Salt,
	}
	json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)
}

func publishEvent(user User) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
