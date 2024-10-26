package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"water-api/middleware"
	"water-api/model"
    "water-api/util"

	"github.com/go-chi/chi/v5"
)



func AuthRoutes(r chi.Router, db *sql.DB) {
    r.Post("/register", RegisterHandler(db))
    r.Post("/login", LoginHandler(db))
    r.With(middleware.AuthMiddleware).Get("/protected", ProtectedHandler)
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var user model.User

        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        _, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
        if err != nil {
            log.Printf("Failed to register user: %v", err)
            http.Error(w, "Failed to register user", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        w.Write([]byte("User registered successfully"))
    }
}


func LoginHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var credentials model.User

        if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        var dbPassword string
        err := db.QueryRow("SELECT password FROM users WHERE username = ?", credentials.Username).Scan(&dbPassword)
        if err == sql.ErrNoRows {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        } else if err != nil {
            http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
            return
        }

        if credentials.Password != dbPassword {
            http.Error(w, "Incorrect password", http.StatusUnauthorized)
            return
        }

        token, err := util.GenerateToken(credentials.Username)
        if err != nil {
            http.Error(w, "Could not generate token", http.StatusInternalServerError)
            return
        }

        w.Write([]byte(token))

        // w.Write([]byte("Login successful"))
    }
}


func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
    username := r.Context().Value("username").(string)
    w.Write([]byte("Hello, " + username + "! This is a protected route."))
}
