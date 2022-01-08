package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/kazimanzurrashid/consents-api-go/handlers"
	"github.com/kazimanzurrashid/consents-api-go/services"
)

func closeDB(db *sql.DB) {
	_ = db.Close()
}

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := path.Dir(currentFile)

	if os.Getenv("POSTGRES_HOST") == "" {
		envFile := filepath.Join(currentDir, "./../../.env")

		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Load(envFile); err != nil {
				log.Fatalf("env file load error: %v", err)
				return
			}
		}
	}

	pgConnectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	db, err := sql.Open("postgres", pgConnectionString)

	if err != nil {
		log.Fatalf("postgres open error: %v", err)
		return
	}

	defer closeDB(db)

	start := time.Now()

	for db.Ping() != nil {
		if start.After(start.Add(10 * time.Second)) {
			closeDB(db)
			log.Print(
				"failed to connect to postgres server even after retrying " +
					"for 10 seconds")
			os.Exit(1)
			return
		}
	}

	schemaFile := filepath.Join(currentDir, "./../../schema.sql")
	schema, err := os.ReadFile(schemaFile)

	if err != nil {
		closeDB(db)
		log.Printf("schema file read error: %v", err)
		os.Exit(1)
		return
	}

	if _, err := db.Exec(string(schema)); err != nil {
		closeDB(db)
		log.Printf("schema file execute error: %v", err)
		os.Exit(1)
		return
	}

	us := services.NewUser(db)
	es := services.NewEvent(db)
	uh := handlers.NewUser(us)
	eh := handlers.NewEvent(es)

	router := mux.NewRouter()

	router.HandleFunc("/users", uh.Create).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", uh.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/users/{id}", uh.Detail).Methods(http.MethodGet)
	router.HandleFunc("/events", eh.Create).Methods(http.MethodPost)
	router.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(struct {
			Result    string `json:"result"`
			Timestamp string `json:"timestamp"`
		}{
			Result:    "ok",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}).Methods(http.MethodGet)

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	err = http.ListenAndServe(port, gh.LoggingHandler(os.Stdout, router))

	if err != nil {
		closeDB(db)
		log.Printf("server listen error: %v", err)
		os.Exit(1)
	}
}
