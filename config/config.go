package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)  
  
var JwtKey []byte // Ubah nama variabel menjadi JwtKey  
  
// LoadEnv memuat variabel lingkungan dari file .env    
func LoadEnv() {    
	err := godotenv.Load()    
	if err != nil {    
		log.Println("Warning: .env file not found. Using system environment variables...")    
	}    
    
	jwtSecret := os.Getenv("JWT_SECRET")    
	if jwtSecret == "" {    
		log.Fatalf("Error: JWT_SECRET is not set. Please set it in the .env file or as an environment variable.")    
	}    
    
	JwtKey = []byte(jwtSecret) // Menggunakan JwtKey  
	log.Printf("Environment loaded successfully.")    
}    
