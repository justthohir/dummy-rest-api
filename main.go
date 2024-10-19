package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type APIEndpoint struct {
	Method      string
	Path        string
	URL         string
	Description string
}

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"` // Add protocol to the Config struct
	BaseURL  string `json:"base_url"` // Add BaseURL for "/demo/dummy-rest-api/"
}

func generateDummyItem(id int) Item {
	return Item{
		ID:    id,
		Name:  fmt.Sprintf("Item %d", id),
		Value: fmt.Sprintf("Value %d", rand.Intn(100)),
	}
}

func getItems(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling GET /api/items request")
	items := make([]Item, 10)
	for i := 0; i < 10; i++ {
		items[i] = generateDummyItem(i + 1)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling GET /api/items/{id} request")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	item := generateDummyItem(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling POST /api/items request")
	newItem := generateDummyItem(rand.Intn(1000) + 100)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newItem)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling PUT /api/items/{id} request")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	updatedItem := Item{
		ID:    id,
		Name:  fmt.Sprintf("Updated Item %d", id),
		Value: fmt.Sprintf("Updated Value %d", rand.Intn(100)),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedItem)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling DELETE /api/items/{id} request")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	response := map[string]string{
		"message": fmt.Sprintf("Item %d has been deleted", id),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func indexHandler(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling GET / request")
		baseURL := config.Protocol + "://" + r.Host + config.BaseURL // Use base URL from config
		endpoints := []APIEndpoint{
			{"GET", "/api/items", baseURL + "api/items", "Get all items"},
			{"GET", "/api/items/{id}", baseURL + "api/items/1", "Get a specific item"},
			{"POST", "/api/items", baseURL + "api/items", "Create a new item"},
			{"PUT", "/api/items/{id}", baseURL + "api/items/1", "Update an item"},
			{"DELETE", "/api/items/{id}", baseURL + "api/items/1", "Delete an item"},
		}

		tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dummy REST API Index</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; padding: 20px; }
        h1 { color: #333; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
        th { background-color: #f2f2f2; }
        tr:nth-child(even) { background-color: #f9f9f9; }
        a { color: #0066cc; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <h1>Dummy REST API Endpoints</h1>
    <table>
        <tr>
            <th>Method</th>
            <th>Path</th>
            <th>URL</th>
            <th>Description</th>
        </tr>
        {{range .}}
        <tr>
            <td>{{.Method}}</td>
            <td>{{.Path}}</td>
            <td><a href="{{.URL}}" target="_blank">{{.URL}}</a></td>
            <td>{{.Description}}</td>
        </tr>
        {{end}}
    </table>
</body>
</html>
`

		t, err := template.New("index").Parse(tmpl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, endpoints)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func main() {
	log.Println("Starting the application...")

	// Read the config file
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Error opening config file:", err)
	}
	defer configFile.Close()

	var config Config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatal("Error decoding config file:", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler(config)).Methods("GET")
	router.HandleFunc("/api/items", getItems).Methods("GET")
	router.HandleFunc("/api/items/{id}", getItem).Methods("GET")
	router.HandleFunc("/api/items", createItem).Methods("POST")
	router.HandleFunc("/api/items/{id}", updateItem).Methods("PUT")
	router.HandleFunc("/api/items/{id}", deleteItem).Methods("DELETE")

	log.Printf("Server is starting on %s:%s...\n", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(config.Host+":"+config.Port, router))
}
