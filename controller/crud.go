package controller

import (
	"akhir-belakang/config"
	"akhir-belakang/helper"
	"akhir-belakang/model"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func GetDataLocation(w http.ResponseWriter, r *http.Request) {
	log.Println("GetDataLocation called")

	// Periksa metode HTTP
	if r.Method != http.MethodGet {
		log.Println("Invalid method:", r.Method)
		helper.WriteResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"message": "Method not allowed",
		})
		return
	}

	// Ambil token dari header Authorization
	tokenString, err := helper.GetTokenFromHeader(r)
	if err != nil {
		log.Printf("Token error: %v", err)
		helper.WriteResponse(w, http.StatusUnauthorized, map[string]string{
			"message": "Unauthorized: " + err.Error(),
		})
		return
	}

	// Periksa apakah token ada di blacklist
	if helper.IsTokenBlacklisted(tokenString) {
		log.Printf("Token is blacklisted: %v", tokenString)
		helper.WriteResponse(w, http.StatusUnauthorized, map[string]string{
			"message": "Unauthorized: Token has been blacklisted",
		})
		return
	}

	// Verifikasi token
	claims := &model.Claims{}
	if err := helper.ParseAndValidateToken(tokenString, claims); err != nil {
		log.Printf("Token validation error: %v", err)
		helper.WriteResponse(w, http.StatusUnauthorized, map[string]string{
			"message": "Unauthorized: " + err.Error(),
		})
		return
	}

	// Ambil data lokasi dari database PostgreSQL
	var lokasi []model.Lokasi
	if err := config.DB.Find(&lokasi).Error; err != nil {
		log.Printf("Database error while fetching locations: %v", err)
		helper.WriteResponse(w, http.StatusInternalServerError, map[string]string{
			"message": "Failed to fetch locations",
		})
		return
	}

	// Jika tidak ada data ditemukan
	if len(lokasi) == 0 {
		log.Println("No locations found in the database")
		helper.WriteResponse(w, http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "No locations found",
			"data":    []model.Lokasi{},
		})
		return
	}

	// Kirim data lokasi dalam format JSON
	log.Printf("Locations retrieved successfully: %d records", len(lokasi))
	helper.WriteResponse(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Locations retrieved successfully",
		"data":    lokasi,
	})
}

func CreateDataLocations(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateDataLocations called")

	// Periksa metode HTTP
	if r.Method != http.MethodPost {
		log.Println("Invalid method:", r.Method)
		helper.WriteResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"message": "Method not allowed",
		})
		return
	}

	// Ambil token dari header Authorization
	tokenString, err := helper.GetTokenFromHeader(r)
	if err != nil {
		log.Printf("Token error: %v", err)
		helper.WriteResponse(w, http.StatusUnauthorized, map[string]string{
			"message": "Unauthorized: " + err.Error(),
		})
		return
	}

	// Periksa apakah token ada di blacklist
	if helper.IsTokenBlacklisted(tokenString) {
		log.Printf("Token is blacklisted: %v", tokenString)
		helper.WriteResponse(w, http.StatusUnauthorized, map[string]string{
			"message": "Unauthorized: Token has been blacklisted",
		})
		return
	}

	// Verifikasi token
	claims := &model.Claims{}
	if err := helper.ParseAndValidateToken(tokenString, claims); err != nil {
		log.Printf("Token validation error: %v", err)
		helper.WriteResponse(w, http.StatusUnauthorized, map[string]string{
			"message": "Unauthorized: " + err.Error(),
		})
		return
	}

	// Decode body request
	var lokasi model.Lokasi
	if err := json.NewDecoder(r.Body).Decode(&lokasi); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		helper.WriteResponse(w, http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
		return
	}

	// Validasi data
	if lokasi.Nama == "" || lokasi.Alamat == "" || lokasi.Deskripsi == "" {
		log.Println("Validation error: Missing required fields")
		helper.WriteResponse(w, http.StatusBadRequest, map[string]string{
			"message": "All fields (nama, alamat, deskripsi) are required",
		})
		return
	}

	// Atur waktu saat ini ke CreatedAt
	lokasi.CreatedAt = time.Now()

	// Simpan ke database
	if err := config.DB.Create(&lokasi).Error; err != nil {
		log.Printf("Failed to save location to database: %v", err)
		helper.WriteResponse(w, http.StatusInternalServerError, map[string]string{
			"message": "Failed to create location",
		})
		return
	}

	// Kirim respons sukses
	log.Printf("Location created successfully: %+v", lokasi)
	helper.WriteResponse(w, http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "Location created successfully",
		"data":    lokasi,
	})
}
