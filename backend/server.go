package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func health(w http.ResponseWriter, r *http.Request) {
	slog.Info("Health check", "method", r.Method, "path", r.URL.Path)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s\n", port)
	err := http.ListenAndServe(":"+port, mux)
	if err == http.ErrServerClosed {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
