package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Image struct {
	ID   int    `json:"id"`
	Link string `json:"link"`
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Open a connection to the database
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_DATABASE"),
	))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new HTTP router using gorilla/mux
	router := mux.NewRouter()

	// Register the GET /images endpoint with a handler function
	router.HandleFunc("/images", func(w http.ResponseWriter, r *http.Request) {
		// Specify the allowed origin of your React frontend
		allowedOrigin := "http://localhost:3000" // Change this to your React app's origin

		// Check the origin of the incoming request
		origin := r.Header.Get("Origin")
		if origin != allowedOrigin {
			http.Error(w, "Not allowed", http.StatusForbidden)
			return
		}

		// Add CORS headers to allow requests from your React frontend
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Check for preflight requests (OPTIONS) and handle them
		if r.Method == "OPTIONS" {
			return
		}

		// Get all of the image links from the database
		rows, err := db.Query("SELECT id, link FROM images")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Create a slice to store the image links
		var images []Image
		for rows.Next() {
			var image Image
			err := rows.Scan(&image.ID, &image.Link)
			if err != nil {
				log.Fatal(err)
			}

			images = append(images, image)
		}

		// Encode the image links to JSON and write them to the response
		json, err := json.Marshal(images)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	})

	// Attach the router to the HTTP server
	http.Handle("/", router)

	// Start the HTTP server
	fmt.Println("Server is running on :8000")
	http.ListenAndServe(":8000", nil)
}
