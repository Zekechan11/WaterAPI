package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"water-api/api"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Printf("DB ERR: %v", err)
		log.Fatal("Failed to connect to the database:", err)
	}

	r := chi.NewRouter()

	api.AuthRoutes(r, db)

	if err := http.ListenAndServe(":3030", r); err != nil {
		log.Fatal("Error starting server:", err)
	}

	fmt.Println("Server starting on http://localhost:3030")
}
