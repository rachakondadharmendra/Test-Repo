package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"backend_golang/src/db"
	"backend_golang/src/handlers"
	"backend_golang/src/logger"
)

type APICall struct {
	Method      string
	Path        string
	Description string
	Usage       string
	Headers     map[string]string
	Body        string
	Example     string
}

func AllAPIHandler(w http.ResponseWriter, r *http.Request) {
	apiCalls := []APICall{
		{
			Method:      "POST",
			Path:        "/api/insertdata",
			Description: "Inserts new data into the database.",
			Usage:       "Send a POST request to /api/insertdata.",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{
    "name": "John Doe",
    "email": "john@example.com",
    "message": "Hello, World!",
    "status": true
}`,
			Example: `curl -X POST -H 'Content-Type: application/json' \
-d '{"name": "John Doe", "email": "john@example.com", "message": "Hello, World!", "status": true}' \
http://localhost:8080/api/insertdata`,
		},
		{
			Method:      "GET",
			Path:        "/api/getdata",
			Description: "Retrieves data from the database.",
			Usage:       "Send a GET request to /api/getdata.",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Example: `curl http://localhost:8080/api/getdata`,
		},
		{
			Method:      "PUT",
			Path:        "/api/updatedata/{id}",
			Description: "Updates existing data in the database.",
			Usage:       "Send a PUT request to /api/updatedata/{id}.",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{
    "id": "12345678",
    "name": "Updated Name",
    "email": "updated@example.com",
    "message": "Updated message",
    "status": false
}`,
			Example: `curl -X PUT -H 'Content-Type: application/json' \
-d '{"id": "12345678", "name": "Updated Name", "email": "updated@example.com", "message": "Updated message", "status": false}' \
http://localhost:8080/api/updatedata/12345678`,
		},
		{
			Method:      "DELETE",
			Path:        "/api/deletedata/{id}",
			Description: "Deletes data from the database.",
			Usage:       "Send a DELETE request to /api/deletedata/{id}.",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Example: `curl -X DELETE http://localhost:8080/api/deletedata/12345678`,
		},
		{
			Method:      "PATCH",
			Path:        "/api/patchdata/{id}",
			Description: "Partially updates existing data in the database.",
			Usage:       "Send a PATCH request to /api/patchdata/{id}.",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{
    "status": false
}`,
			Example: `curl -X PATCH -H 'Content-Type: application/json' \
-d '{"status": false}' \
http://localhost:8080/api/patchdata/12345678`,
		},
	}

	// Define a custom template function to convert a string to lowercase
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	// Create a new template and parse the HTML file
	tmpl, err := template.New("api_all.html").Funcs(funcMap).ParseFiles("api_all.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the API calls data
	err = tmpl.Execute(w, apiCalls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    // Render the error page template
    tmpl, err := template.ParseFiles("api_error.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Execute the template
    err = tmpl.Execute(w, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
	// Initialize logger
	logger := logger.InitLogger()

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file:", err)
	}

	// Connect to MongoDB
	err = db.ConnectMongoDB(logger)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.DisconnectMongoDB(logger)

	// Define HTTP routes with CORS configuration
	router := mux.NewRouter()
	router.HandleFunc("/api/insertdata", handlers.InsertDataHandler).Methods("POST")
	router.HandleFunc("/api/getdata", handlers.GetDataHandler).Methods("GET")
	router.HandleFunc("/api/updatedata/{id}", handlers.UpdateDataHandler).Methods("PUT")
	router.HandleFunc("/api/deletedata/{id}", handlers.DeleteDataHandler).Methods("DELETE")
	router.HandleFunc("/api/patchdata/{id}", handlers.PatchDataHandler).Methods("PATCH")

	// Add the /api/all route
	router.HandleFunc("/api/all", AllAPIHandler).Methods("GET")

	// Configure CORS options
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"}, // Replace with your frontend domain(s)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}

	// Apply CORS middleware
	corsHandler := cors.New(corsOptions).Handler(router)

	// Handle invalid URLs
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	// Start HTTP server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("Server is running on port %s...\n", port)
	logger.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(addr, corsHandler))
}
