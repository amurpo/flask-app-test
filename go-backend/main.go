package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Image struct to map the database records to JSON format
type Image struct {
	ID   int    `json:"id"`
	Link string `json:"link"`
}

// Middleware to set CORS headers
func setCORSHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := os.Getenv("REACT_APP_ORIGIN")
		if origin == "" {
			origin = "*" // Default to allow all origins if not set
		}

		// If the origin uses HTTPS, enforce HTTPS
		if strings.HasPrefix(origin, "https://") {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Attempt to load .env file, but don't fail if it's not found
	_ = godotenv.Load()

	// MySQL database credentials from environment variables
	mysqlUser := getEnv("MYSQL_USER", "root")
	mysqlPassword := getEnv("MYSQL_PASSWORD", "")
	mysqlHost := getEnv("MYSQL_HOST", "localhost")
	mysqlDatabase := getEnv("MYSQL_DATABASE", "test")

	// Open a connection to the database
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s",
		mysqlUser,
		mysqlPassword,
		mysqlHost,
		mysqlDatabase,
	))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new HTTP router using gorilla/mux
	router := mux.NewRouter()

	// Register the GET /images endpoint with a handler function
	router.HandleFunc("/images", func(w http.ResponseWriter, r *http.Request) {
		// Query to get all image links from the database
		rows, err := db.Query("SELECT id, link FROM images")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Create a slice to store the image links
		var images []Image
		for rows.Next() {
			var image Image
			if err := rows.Scan(&image.ID, &image.Link); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			images = append(images, image)
		}

		// Encode the image links to JSON and write them to the response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(images); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("GET", "OPTIONS")

	// Apply the CORS middleware to the router
	http.Handle("/", setCORSHeaders(router))

	// Start the HTTP server
	port := getEnv("PORT", "8000")
	fmt.Printf("Server is running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// getEnv retrieves the value of the environment variable named by the key.
// If the variable is not present, it returns the fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
