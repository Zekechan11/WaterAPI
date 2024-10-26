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
	r.With(middleware.AuthMiddleware).Group(func(r chi.Router) {
        r.Post("/add", AddDataHandler(db))
        r.Get("/get", GetDataHandler(db))
        r.Patch("/update", UpdateDataHandler(db))
        r.Delete("/delete", DeleteDataHandler(db))
    })
}

// Get Data
// http://localhost:3030/add
// body:
// {
// 	"name": "name"
// }
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
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": "Bottle added successfully"}
		json.NewEncoder(w).Encode(response)
	}
}

// Get Data
// http://localhost:3030/get
func GetDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT * FROM bottles")
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Failed to to get bottle", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var bottles []model.Bottle
		for rows.Next() {
			var bottle model.Bottle
			if err := rows.Scan(&bottle.Id, &bottle.Name); err != nil {
				http.Error(w, "Failed to scan bottle", http.StatusInternalServerError)
				return
			}
			bottles = append(bottles, bottle)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Row error: %v", err)
			http.Error(w, "Failed to get bottles", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(bottles); err != nil {
			http.Error(w, "Failed to encode bottles to JSON", http.StatusInternalServerError)
		}
	}
}

// Update Data
// http://localhost:3030/update?id={id}
// body:
// {
// 	"name": "new name"
// }
func UpdateDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		var bottle model.Bottle
		if err := json.NewDecoder(r.Body).Decode(&bottle); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		stmt, err := db.Prepare("UPDATE bottles SET name = ? WHERE id = ?")
		if err != nil {
			http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(bottle.Name, id)
		if err != nil {
			http.Error(w, "Failed to update data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"message": "Data updated successfully"}
		json.NewEncoder(w).Encode(response)
	}
}

// Delete Data
// http://localhost:3030/delete?id={id}
func DeleteDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		_, err := db.Exec("DELETE FROM bottles WHERE id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		res := map[string]string{"message": "Bottle Deleted!"}
		json.NewEncoder(w).Encode(res)
	}
}
