package main

import (
    "bytes"
    "io"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
)

func enableCORS(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
}

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

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading request body", http.StatusBadRequest)
        return
    }

    oauthTokenURL := os.Getenv("OAUTH_TOKEN_URL")
    proxyReq, err := http.NewRequest(http.MethodPost, oauthTokenURL, bytes.NewBuffer(body))
    if err != nil {
        http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
        return
    }

    proxyReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    proxyReq.Header.Set("Accept", "application/json")

    client := &http.Client{}
    resp, err := client.Do(proxyReq)
    if err != nil {
        http.Error(w, "Error forwarding request", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }
    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}

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

    accessToken := r.URL.Query().Get("access_token")
    if accessToken == "" {
        http.Error(w, "Missing access token", http.StatusBadRequest)
        return
    }

    userInfoURL := os.Getenv("USER_INFO_URL")
    proxyReq, err := http.NewRequest(http.MethodGet, userInfoURL, nil)
    if err != nil {
        http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
        return
    }

    proxyReq.Header.Set("Authorization", "Bearer "+accessToken)

    client := &http.Client{}
    resp, err := client.Do(proxyReq)
    if err != nil {
        http.Error(w, "Error forwarding request", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }
    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found, falling back to environment variables")
    }

    http.HandleFunc("/proxy/oauth/token", oauthTokenProxy)
    http.HandleFunc("/proxy/user-info", userInfoProxy)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
