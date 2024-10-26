package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"water-api/middleware"
	"water-api/model"

	"github.com/go-chi/chi/v5"
)

func AdminRoutes(r chi.Router, db *sql.DB) {
	r.With(middleware.AuthMiddleware).
		Post("/adddata", AddDataHandler(db))
}

func AddDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bottle model.Bottle

        if err := json.NewDecoder(r.Body).Decode(&bottle); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }
		_, err := db.Exec("INSERT INTO bottles (name) VALUES (?)", bottle.Name)
        if err != nil {
			log.Printf("%v", err)
            http.Error(w, "Failed to to add bottle", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        w.Write([]byte("Bottle added"))
	}
}

// func GetDataHandler(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		data, err := db.Exec("SELECT * FROM bottles")
//         if err != nil {
// 			log.Printf("%v", err)
//             http.Error(w, "Failed to to get bottle", http.StatusInternalServerError)
//             return
//         }

//         w.Write([]byte(data))
// 	}
// }

