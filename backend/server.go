package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Yongbeom-Kim/transfer/backend/internal/middleware"
)

func health(w http.ResponseWriter, r *http.Request) {
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
	err := http.ListenAndServe(":"+port,
		middleware.Compose(
			middleware.CORSMiddleware,
			middleware.Logger,
		)(mux),
	)
	if err == http.ErrServerClosed {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
