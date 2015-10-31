package server

import (
	"log"
	"net/http"
	"os"
)

// Start is a convenience method for starting the server, listening on PORT
func Start() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not set")
	}
	log.Printf("Starting HTTP server on %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server start error: ", err)
	}
}
