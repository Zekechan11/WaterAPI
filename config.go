package main

import (
    "database/sql"
    "fmt"
)

func InitDB(db *sql.DB) {
    _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT,
        password TEXT
    )`)
    if err != nil {
        fmt.Println("Error creating users table:", err)
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS bottles (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT
    )`)
    if err != nil {
        fmt.Println("Error creating bottles table:", err)
    }
}
