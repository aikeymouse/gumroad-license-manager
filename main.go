package main

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Config struct {
	GumroadToken string `json:"gumroad_token"`
}

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type ProductsResponse struct {
	Success  bool      `json:"success"`
	Products []Product `json:"products"`
}

type License struct {
	ID             string `json:"id"`
	ProductName    string `json:"product_name"`
	LicenseKey     string `json:"license_key"`
	Permalink      string `json:"permalink"`
	SaleDatetime   string `json:"sale_datetime"`
	PurchaserEmail string `json:"purchaser_email"`
	Refunded       bool   `json:"refunded"`
	Disputed       bool   `json:"disputed"`
	Chargebacked   bool   `json:"chargebacked"`
}

type LicensesResponse struct {
	Success  bool      `json:"success"`
	Licenses []License `json:"licenses"`
}

type Sale struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	Price           int    `json:"price"`
	GumroadFee      int    `json:"gumroad_fee"`
	Currency        string `json:"currency"`
	Quantity        int    `json:"quantity"`
	DiscoverFee     int    `json:"discover_fee"`
	CanContact      bool   `json:"can_contact"`
	Referrer        string `json:"referrer"`
	OrderID         int64  `json:"order_id"`
	CreatedAt       string `json:"created_at"`
	ProductID       string `json:"product_id"`
	ProductName     string `json:"product_name"`
	Refunded        bool   `json:"refunded"`
	Disputed        bool   `json:"disputed"`
	Chargebacked    bool   `json:"chargebacked"`
	AffiliateCredit int    `json:"affiliate_credit"`
	// Adding some common fields from API response
	PurchaserID string `json:"purchaser_id"`
	LicenseKey  string `json:"license_key"`
	Timestamp   string `json:"timestamp"`
	Daystamp    string `json:"daystamp"`
}

type SalesResponse struct {
	Success bool   `json:"success"`
	Sales   []Sale `json:"sales"`
}

type ValidateLicenseRequest struct {
	ProductID  string `json:"product_id"`
	LicenseKey string `json:"license_key"`
}

type LicenseValidationResponse struct {
	Success  bool                   `json:"success"`
	Uses     int                    `json:"uses,omitempty"`
	Purchase map[string]interface{} `json:"purchase,omitempty"`
	Message  string                 `json:"message,omitempty"`
}

type APICall struct {
	Timestamp    time.Time
	Method       string
	URL          string
	Status       int
	Duration     time.Duration
	Error        string
	RequestBody  string
	ResponseBody string
	Headers      map[string]string
}

type PageData struct {
	Title          string
	CurrentPage    string
	BackLink       string
	Products       []Product
	Licenses       []License
	Sales          []Sale
	ProductID      string
	APICallsResult []APICall
}

type App struct {
	config    Config
	apiCalls  []APICall
	mu        sync.RWMutex
	templates *template.Template
}

func loadConfig() (Config, error) {
	var config Config
	file, err := os.Open("config.json")
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func saveConfig(config Config) error {
	file, err := os.Create("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

func isTokenConfigured(config Config) bool {
	return config.GumroadToken != "" && config.GumroadToken != "YOUR_GUMROAD_ACCESS_TOKEN_HERE"
}

func (app *App) loadTemplates() error {
	funcMap := template.FuncMap{
		"div":      func(a, b float64) float64 { return a / b },
		"mul":      func(a, b int) time.Duration { return time.Duration(a * b) },
		"mulF":     func(a int, b float64) float64 { return float64(a) * b },
		"eq":       func(a, b string) bool { return a == b },
		"eqInt":    func(a, b int) bool { return a == b },
		"unescape": func(s string) template.HTML { return template.HTML(html.UnescapeString(s)) },
		"jsonMarshal": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		},
		"sub": func(a, b int) int { return a - b },
		"durationMs": func(d time.Duration) int {
			return int(d.Nanoseconds() / 1000000)
		},
	}

	templates := template.New("").Funcs(funcMap)

	// Parse all template files
	var err error
	templates, err = templates.ParseGlob("templates/*.html")
	if err != nil {
		return err
	}

	app.templates = templates
	return nil
}

func (app *App) logAPICall(method, url string, status int, duration time.Duration, err error, requestBody, responseBody string, headers map[string]string) {
	app.mu.Lock()
	defer app.mu.Unlock()

	apiCall := APICall{
		Timestamp:    time.Now(),
		Method:       method,
		URL:          url,
		Status:       status,
		Duration:     duration,
		RequestBody:  requestBody,
		ResponseBody: responseBody,
		Headers:      headers,
	}

	if err != nil {
		apiCall.Error = err.Error()
	}

	app.apiCalls = append(app.apiCalls, apiCall)

	// Keep only last 100 API calls
	if len(app.apiCalls) > 100 {
		app.apiCalls = app.apiCalls[len(app.apiCalls)-100:]
	}
}

func (app *App) makeGumroadRequest(url string) ([]byte, error) {
	start := time.Now()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		app.logAPICall("GET", url, 0, time.Since(start), err, "", "", nil)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+app.config.GumroadToken)

	// Capture request headers
	headers := make(map[string]string)
	for k, v := range req.Header {
		headers[k] = v[0]
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		app.logAPICall("GET", url, 0, time.Since(start), err, "", "", headers)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	responseBody := string(body)

	app.logAPICall("GET", url, resp.StatusCode, time.Since(start), err, "", responseBody, headers)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, err
}

func (app *App) getProducts() ([]Product, error) {
	body, err := app.makeGumroadRequest("https://api.gumroad.com/v2/products")
	if err != nil {
		return nil, err
	}

	var response ProductsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("API request was not successful")
	}

	return response.Products, nil
}

func (app *App) getLicenses(productID string) ([]License, error) {
	url := fmt.Sprintf("https://api.gumroad.com/v2/products/%s/subscribers", productID)
	body, err := app.makeGumroadRequest(url)
	if err != nil {
		return nil, err
	}

	var response LicensesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("API request was not successful")
	}

	return response.Licenses, nil
}

func (app *App) getSales(productID string) ([]Sale, error) {
	url := fmt.Sprintf("https://api.gumroad.com/v2/sales?product_id=%s", productID)

	body, err := app.makeGumroadRequest(url)
	if err != nil {
		return nil, err
	}

	var response SalesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("API request was not successful")
	}

	return response.Sales, nil
}

func (app *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Index handler called")
	products, err := app.getProducts()
	if err != nil {
		log.Printf("Failed to fetch products: %v", err)
		http.Error(w, "Failed to fetch products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Fetched %d products successfully", len(products))
	data := PageData{
		Title:       "Products",
		CurrentPage: "products",
		Products:    products,
	}

	w.Header().Set("Content-Type", "text/html")
	err = app.templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Products template rendered successfully")
}

func (app *App) licensesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	indexStr := vars["index"]

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "Invalid product index", http.StatusBadRequest)
		return
	}

	// Get products to find the actual product ID from the index
	products, err := app.getProducts()
	if err != nil {
		http.Error(w, "Failed to fetch products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if index < 0 || index >= len(products) {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	product := products[index]
	productID := product.ID

	licenses, err := app.getLicenses(productID)
	if err != nil {
		http.Error(w, "Failed to fetch licenses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:       fmt.Sprintf("License Keys - %s", product.Name),
		CurrentPage: "licenses",
		BackLink:    "/",
		Licenses:    licenses,
		ProductID:   productID,
	}

	w.Header().Set("Content-Type", "text/html")
	err = app.templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (app *App) salesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	indexStr := vars["index"]

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "Invalid product index", http.StatusBadRequest)
		return
	}

	// Get products to find the actual product ID from the index
	products, err := app.getProducts()
	if err != nil {
		http.Error(w, "Failed to fetch products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if index < 0 || index >= len(products) {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	product := products[index]
	productID := product.ID

	sales, err := app.getSales(productID)
	if err != nil {
		http.Error(w, "Failed to fetch sales: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:       fmt.Sprintf("Sales - %s", product.Name),
		CurrentPage: "sales",
		BackLink:    "/",
		Sales:       sales,
		ProductID:   productID,
	}

	w.Header().Set("Content-Type", "text/html")
	err = app.templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (app *App) apiLogHandler(w http.ResponseWriter, r *http.Request) {
	app.mu.RLock()
	apiCalls := make([]APICall, len(app.apiCalls))
	copy(apiCalls, app.apiCalls)
	app.mu.RUnlock()

	// Reverse the slice to show newest first
	for i, j := 0, len(apiCalls)-1; i < j; i, j = i+1, j-1 {
		apiCalls[i], apiCalls[j] = apiCalls[j], apiCalls[i]
	}

	// Determine back link based on referer
	backLink := "/"
	referer := r.Header.Get("Referer")
	if referer != "" {
		// Parse the referer URL to get the path
		if refererURL, err := url.Parse(referer); err == nil {
			refererPath := refererURL.Path
			// Only use referer if it's from our application and not the same page
			if refererPath != "/api-log" && (refererPath == "/" ||
				strings.HasPrefix(refererPath, "/licenses/") ||
				strings.HasPrefix(refererPath, "/sales/")) {
				backLink = refererPath
			}
		}
	}

	data := PageData{
		Title:          "API Call Log",
		CurrentPage:    "api-log",
		BackLink:       backLink,
		APICallsResult: apiCalls,
	}

	w.Header().Set("Content-Type", "text/html")
	err := app.templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (app *App) setupHandler(w http.ResponseWriter, r *http.Request) {
	// Check if token is already configured
	config, err := loadConfig()
	log.Printf("Setup handler: token=%s, err=%v", config.GumroadToken, err)
	if err == nil && config.GumroadToken != "" && config.GumroadToken != "YOUR_GUMROAD_ACCESS_TOKEN_HERE" {
		// Token is configured, redirect to main page
		log.Printf("Token is configured, redirecting to /")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	log.Printf("Showing setup page")
	data := PageData{
		Title:       "Setup - Gumroad Token",
		CurrentPage: "setup",
	}

	w.Header().Set("Content-Type", "text/html")
	err = app.templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (app *App) setupSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid JSON data",
		})
		return
	}

	if requestData.Token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Token cannot be empty",
		})
		return
	}

	// Test the token by making a simple API call
	if err := app.testGumroadToken(requestData.Token); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid token: " + err.Error(),
		})
		return
	}

	// Save the token
	app.config.GumroadToken = requestData.Token
	if err := saveConfig(app.config); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to save token: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Token saved successfully",
	})
}

func (app *App) testGumroadToken(token string) error {
	req, err := http.NewRequest("GET", "https://api.gumroad.com/v2/products", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return fmt.Errorf("unauthorized - invalid token")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return nil
}

func (app *App) setupMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check current token status dynamically
		config, err := loadConfig()
		log.Printf("Setup middleware: path=%s, token=%s, err=%v", r.URL.Path, config.GumroadToken, err)
		if err != nil || config.GumroadToken == "" || config.GumroadToken == "YOUR_GUMROAD_ACCESS_TOKEN_HERE" {
			// No token configured or placeholder token, redirect to setup
			log.Printf("Redirecting to setup: no valid token")
			http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
			return
		}
		log.Printf("Token is valid, continuing to handler")
		next.ServeHTTP(w, r)
	})
}

// apiCallsJSONHandler returns API calls data as JSON
func (app *App) apiCallsJSONHandler(w http.ResponseWriter, r *http.Request) {
	app.mu.RLock()
	apiCallsCopy := make([]APICall, len(app.apiCalls))
	copy(apiCallsCopy, app.apiCalls)
	app.mu.RUnlock()

	// Reverse the slice to show newest first (same as template)
	for i, j := 0, len(apiCallsCopy)-1; i < j; i, j = i+1, j-1 {
		apiCallsCopy[i], apiCallsCopy[j] = apiCallsCopy[j], apiCallsCopy[i]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiCallsCopy)
}

func (app *App) validateLicenseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ValidateLicenseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.ProductID == "" || req.LicenseKey == "" {
		http.Error(w, "Missing product_id or license_key", http.StatusBadRequest)
		return
	}

	// Call Gumroad license verification API
	url := "https://api.gumroad.com/v2/licenses/verify"

	// Create form data for the request with increment_uses_count=false
	data := fmt.Sprintf("product_id=%s&license_key=%s&increment_uses_count=false", req.ProductID, req.LicenseKey)

	startTime := time.Now()
	httpReq, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		duration := time.Since(startTime)
		app.logAPICall("POST", url, 0, duration, err, data, "", nil)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		duration := time.Since(startTime)
		app.logAPICall("POST", url, 0, duration, err, data, "", nil)
		http.Error(w, "Failed to validate license", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	responseBody := string(body)
	duration := time.Since(startTime)

	// Log the API call
	headers := make(map[string]string)
	headers["Content-Type"] = httpReq.Header.Get("Content-Type")
	app.logAPICall("POST", url, resp.StatusCode, duration, err, data, responseBody, headers)

	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Parse the response
	var gumroadResponse map[string]interface{}
	err = json.Unmarshal(body, &gumroadResponse)
	if err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	// Create our response
	response := LicenseValidationResponse{
		Success: false,
	}

	if success, ok := gumroadResponse["success"].(bool); ok && success {
		response.Success = true

		if uses, ok := gumroadResponse["uses"].(float64); ok {
			response.Uses = int(uses)
		}

		if purchase, ok := gumroadResponse["purchase"].(map[string]interface{}); ok {
			response.Purchase = purchase
		}
	} else {
		if msg, ok := gumroadResponse["message"].(string); ok {
			response.Message = msg
		} else {
			response.Message = "Invalid license key"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	app := &App{
		config:   config,
		apiCalls: make([]APICall, 0),
	}

	// Load templates
	err = app.loadTemplates()
	if err != nil {
		log.Fatal("Failed to load templates:", err)
	}

	r := mux.NewRouter()

	// Setup routes (always available)
	r.HandleFunc("/setup", app.setupHandler).Methods("GET")
	r.HandleFunc("/setup/submit", app.setupSubmitHandler).Methods("POST")

	// Favicon handler (returns empty response)
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}).Methods("GET")

	// Main application routes with setup middleware
	r.HandleFunc("/", app.setupMiddleware(app.indexHandler)).Methods("GET")
	r.HandleFunc("/licenses/{index:[0-9]+}", app.setupMiddleware(app.licensesHandler)).Methods("GET")
	r.HandleFunc("/sales/{index:[0-9]+}", app.setupMiddleware(app.salesHandler)).Methods("GET")
	r.HandleFunc("/api-log", app.setupMiddleware(app.apiLogHandler)).Methods("GET")
	r.HandleFunc("/api/api-calls", app.setupMiddleware(app.apiCallsJSONHandler)).Methods("GET")
	r.HandleFunc("/validate-license", app.setupMiddleware(app.validateLicenseHandler)).Methods("POST")

	// Static file server (always available)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	log.Printf("Server starting on port %s", port)
	if !isTokenConfigured(config) {
		log.Printf("Visit http://localhost:%s/setup to configure your Gumroad token", port)
	} else {
		log.Printf("Visit http://localhost:%s to access the application", port)
	}

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
