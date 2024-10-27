package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

var productStore = make(map[string]Product)

func main() {
	http.HandleFunc("/products", createProductHandler)  // POST /products
	http.HandleFunc("/products/", productHandlerWithID) // GET /products/{id}

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// createProductHandler handles POST requests to create a new product without an ID in the URL
func createProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var product Product

	// Parse the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	product.ID = generateUniqueID()
	productStore[product.ID] = product

	log.Printf("Created Product - ID: %s, Name: %s, Count: %d", product.ID, product.Name, product.Count)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// productHandlerWithID handles GET and POST requests to retrieve or update a product by ID
func productHandlerWithID(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.RequestURI)

	productID := strings.TrimPrefix(r.URL.Path, "/products/")
	if productID == "" {
		http.Error(w, "Product ID not provided", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetProduct(w, productID)
	case http.MethodPost:
		handlePostProduct(w, r, productID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetProduct processes GET requests to retrieve a product by ID
func handleGetProduct(w http.ResponseWriter, productID string) {
	// Check if the product exists in the productStore
	product, exists := productStore[productID]
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Encode and send product as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode product data", http.StatusInternalServerError)
	}
}

// handlePostProduct processes POST requests to update a product by ID
func handlePostProduct(w http.ResponseWriter, r *http.Request, productID string) {
	var product Product
	product.ID = productID // Set the ID from the URL

	// Parse the JSON request body to retrieve name and count
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	productStore[productID] = product

	// log.Printf("Updated Product - ID: %s, Name: %s, Count: %d", product.ID, product.Name, product.Count)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// generateUniqueID generates a random unique ID for each product
func generateUniqueID() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(1000000)) // Generate a random number as ID
}
