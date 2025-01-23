
package main

import (
    "bytes"
    "io"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
)

/* 
  enableCORS: Enables Cross-Origin Resource Sharing (CORS) by setting appropriate headers.
  - Allows requests from any origin (`*`).
  - Supports HTTP methods: POST, OPTIONS, GET.
  - Specifies allowed headers: Content-Type, Accept, Authorization.
*/
func enableCORS(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
}

/* 
  oauthTokenProxy: Proxies OAuth token requests to the external OAuth server.
  - Handles POST and OPTIONS methods.
  - Reads the request body and forwards it to the OAuth token endpoint.
  - Returns the response from the external server to the client.
  - Reads the OAuth token endpoint URL from the environment variable `OAUTH_TOKEN_URL`.
*/
func oauthTokenProxy(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodOptions {
        enableCORS(w)
        w.WriteHeader(http.StatusOK)
        return
    }

    if r.Method != http.MethodPost {
        enableCORS(w)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    enableCORS(w)

    // Read the request body
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading request body", http.StatusBadRequest)
        return
    }

    // Create a new proxy request to the OAuth token endpoint
    oauthTokenURL := os.Getenv("OAUTH_TOKEN_URL")
    proxyReq, err := http.NewRequest(http.MethodPost, oauthTokenURL, bytes.NewBuffer(body))
    if err != nil {
        http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
        return
    }

    // Set headers for the proxy request
    proxyReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    proxyReq.Header.Set("Accept", "application/json")

    // Forward the request to the external OAuth server
    client := &http.Client{}
    resp, err := client.Do(proxyReq)
    if err != nil {
        http.Error(w, "Error forwarding request", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    // Copy response headers and body from the external server to the client
    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }
    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}

/* 
  userInfoProxy: Proxies user info requests to an external user information server.
  - Handles GET and OPTIONS methods.
  - Extracts the `access_token` query parameter from the client request.
  - Forwards the request to the user info endpoint with the `Authorization` header.
  - Returns the response from the external server to the client.
  - Reads the user info endpoint URL from the environment variable `USER_INFO_URL`.
*/
func userInfoProxy(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodOptions {
        enableCORS(w)
        w.WriteHeader(http.StatusOK)
        return
    }

    if r.Method != http.MethodGet {
        enableCORS(w)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    enableCORS(w)

    // Get the access token from query parameters
    accessToken := r.URL.Query().Get("access_token")
    if accessToken == "" {
        http.Error(w, "Missing access token", http.StatusBadRequest)
        return
    }

    // Create a new proxy request to the user info endpoint
    userInfoURL := os.Getenv("USER_INFO_URL")
    proxyReq, err := http.NewRequest(http.MethodGet, userInfoURL, nil)
    if err != nil {
        http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
        return
    }

    // Set the Authorization header with the access token
    proxyReq.Header.Set("Authorization", "Bearer "+accessToken)

    // Forward the request to the external user info server
    client := &http.Client{}
    resp, err := client.Do(proxyReq)
    if err != nil {
        http.Error(w, "Error forwarding request", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    // Copy response headers and body from the external server to the client
    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }
    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}

/* 
  main: Entry point of the application.
  - Loads environment variables from a `.env` file (if present).
  - Registers HTTP handlers for the OAuth token proxy and user info proxy.
  - Starts the HTTP server on port 8080.
*/
func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found, falling back to environment variables")
    }

    // Register HTTP handlers
    http.HandleFunc("/proxy/oauth/token", oauthTokenProxy)
    http.HandleFunc("/proxy/user-info", userInfoProxy)

    // Start the HTTP server
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
