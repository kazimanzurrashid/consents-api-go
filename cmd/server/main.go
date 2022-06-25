package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/kazimanzurrashid/consents-api-go/handlers"
	"github.com/kazimanzurrashid/consents-api-go/services"
)

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

	closeDB := func() {
		_ = db.Close()
	}

	defer closeDB()

	schemaFile := filepath.Join(currentDir, "./../../schema.sql")
	schema, err := os.ReadFile(schemaFile)

	if err != nil {
		closeDB()
		log.Fatalf("schema file read error: %v", err)
		return
	}

	if _, err := db.Exec(string(schema)); err != nil {
		closeDB()
		log.Fatalf("schema file execute error: %v", err)
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

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: gh.LoggingHandler(os.Stdout, router),
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			closeDB()
			log.Fatalf("server listen and serve error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(
		context.Background(),
		10*time.Second)

	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		closeDB()
		log.Fatalf("server shutdown error: %v", err)
	}
}
