package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)  
  
var DB *gorm.DB  
  
// ConnectDatabase menghubungkan aplikasi ke database PostgreSQL  
func ConnectDatabase() {  
    // Memuat variabel dari file .env  
    err := godotenv.Load()  
    if err != nil {  
        log.Fatal("Error loading .env file")  
    }  
  
    // Ambil variabel dari .env  
    dbUser := os.Getenv("POSTGRESUSER")  
    dbPassword := os.Getenv("POSTGRESPASSWORD")  
    dbHost := os.Getenv("POSTGRESHOST")  
    dbPort := os.Getenv("POSTGRESPORT")  
    dbName := os.Getenv("POSTGRESDATABASE")
  
    // Validasi variabel environment  
    if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {  
        log.Fatal("Database configuration is missing in .env file")  
    }  
  
    // Buat DSN  
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",  
        dbHost, dbUser, dbPassword, dbName, dbPort)  
  
    log.Printf("Connecting to database with DSN: %s", dsn)  
  
    // Koneksi ke PostgreSQL  
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})  
    if err != nil {  
        log.Fatalf("Failed to connect to database: %v", err)  
    }  
  
    // Simpan koneksi ke variabel global DB  
    DB = db  
    log.Println("Connected to the database successfully!")
	
}  
