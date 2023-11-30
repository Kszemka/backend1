package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

const (
	host     = "database"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "task_db"
)

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	// Check if the connection to the database is successful
	err = db.Ping()
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err)
	} else {
		log.Println("Connected to the database")
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := getTasks(db)
		if err != nil {
			log.Printf("Error retrieving tasks from the database: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Convert tasks to JSON
		jsonTasks, err := json.Marshal(tasks)
		if err != nil {
			log.Printf("Error converting tasks to JSON: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set Content-Type header to application/json
		w.Header().Set("Content-Type", "application/json")

		// Write JSON response
		w.Write(jsonTasks)
	})

	handler := cors.Default().Handler(mux)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

func getTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT t.name, t.deadline, c.category_name FROM tasks t JOIN task_category c ON t.category_id = c.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Name, &task.Deadline, &task.Category)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

type Task struct {
	Name     string `json:"name"`
	Deadline string `json:"deadline"`
	Category string `json:"category"`
}
