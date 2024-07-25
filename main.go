package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Name string `json:"name"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func addUserHandler(conn *pgx.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		print("started addUserHandler")
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		ctx := context.Background()
		_, err = conn.Exec(ctx, "INSERT INTO chinnu (name) VALUES ($1)", user.Name)
		if err != nil {
			http.Error(w, "Failed to insert data", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "User added successfully")
	}
}

func main() {
	ctx := context.Background()
	// Initialize the database connection
	conn, err := connectDB(ctx)
	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
		return
	}
	defer conn.Close(ctx)
	// Start the HTTP server
	// http.HandleFunc("/", helloHandler)
	http.HandleFunc("/add-user", addUserHandler(conn))

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// connectDB connects to the PostgreSQL database and returns the connection
func connectDB(ctx context.Context) (*pgx.Conn, error) {
	print("1. started")
	connStr := "postgres://postgres:chinnu@localhost:5433/postgres"
	print("2. started")

	if connStr == "" {
		print("3. started")

		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}
	print("4. started")

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return conn, nil
}
