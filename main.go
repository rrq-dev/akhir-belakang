package main

import (
	"akhir-belakang/config"
	"akhir-belakang/routes"
	"log"
	"net/http"
)

func main() {
	// connect to database
	config.ConnectDatabase()

	//load environtment
	config.LoadEnv()

	http.HandleFunc("/", routes.URL)

	//start server
	port := ":8080" 
	log.Printf("Server is running at port %s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}