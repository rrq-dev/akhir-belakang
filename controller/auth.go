package controller

import (
	"akhir-belakang/config"
	"akhir-belakang/model"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {  
	// Ensure the HTTP method is POST  
	if r.Method != http.MethodPost {  
		w.Header().Set("Content-Type", "application/json")  
		w.WriteHeader(http.StatusMethodNotAllowed)  
		json.NewEncoder(w).Encode(map[string]string{  
			"message": "Method not allowed. Please use POST.",  
		})  
		return  
	}  
  
	var requestData model.UserInput  
  
	// Parse JSON request body  
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {  
		log.Printf("Invalid request data: %v", err)  
		w.Header().Set("Content-Type", "application/json")  
		w.WriteHeader(http.StatusBadRequest)  
		json.NewEncoder(w).Encode(map[string]string{  
			"message": "Invalid input. Please check your request data.",  
		})  
		return  
	}  
  
	// Validate input  
	if requestData.Password != requestData.ConfirmPassword {  
		w.Header().Set("Content-Type", "application/json")  
		w.WriteHeader(http.StatusBadRequest)  
		json.NewEncoder(w).Encode(map[string]string{  
			"message": "Passwords do not match.",  
		})  
		return  
	}  
  
	if requestData.Email == "" || requestData.Username == "" || requestData.Password == "" {  
		w.Header().Set("Content-Type", "application/json")  
		w.WriteHeader(http.StatusBadRequest)  
		json.NewEncoder(w).Encode(map[string]string{  
			"message": "All fields are required.",  
		})  
		return  
	}  
  
	// Check if email or username already exists  
	var existingUser model.Users  
	if err := config.DB.Where("email = ? OR username = ?", requestData.Email, requestData.Username).First(&existingUser).Error; err == nil {  
		log.Printf("User already exists with email: %s or username: %s", requestData.Email, requestData.Username)  
		w.Header().Set("Content-Type", "application/json")  
		w.WriteHeader(http.StatusBadRequest)  
		json.NewEncoder(w).Encode(map[string]string{  
			"message": "Email or username already exists. Please use a different one.",  
		})  
		return  
	}  
  
	// Hash password  
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)  
	if err != nil {  
		log.Printf("Failed to hash password: %v", err)  
		w.Header().Set("Content-Type", "application/json")  
		w.WriteHeader(http.StatusInternalServerError)  
		json.NewEncoder(w).Encode(map[string]string{  
			"message": "Failed to hash password.",  
		})  
		return  
	}  
  
	// Save user to database with default role_id = 2  
	user := model.Users{  
		Email:     requestData.Email,  
		Username:  requestData.Username,  
		Password:  string(hashedPassword),  
		RoleID:    2, // Default role for regular users  
		CreatedAt: time.Now(),  
	}  
  
	if err := config.DB.Create(&user).Error; err != nil {  
		log.Printf("Failed to create user: %v", err)  
		w.Header().Set("Content-Type", "application/json")  
		w.WriteHeader(http.StatusInternalServerError)  
		json.NewEncoder(w).Encode(map[string]string{  
			"message": "Failed to register user. Please try again later.",  
		})  
		return  
	}  
  
	// Send success response  
	w.Header().Set("Content-Type", "application/json")  
	w.WriteHeader(http.StatusCreated)  
	json.NewEncoder(w).Encode(map[string]interface{}{  
		"message": "User registered successfully",  
		"user": map[string]interface{}{  
			"id":       user.ID,  
			"email":    user.Email,  
			"username": user.Username,  
		},  
	})  
}

// Login handles user login    
func Login(w http.ResponseWriter, r *http.Request) {    
	// Ensure the HTTP method is POST    
	if r.Method != http.MethodPost {    
		w.Header().Set("Content-Type", "application/json")    
		w.WriteHeader(http.StatusMethodNotAllowed)    
		json.NewEncoder(w).Encode(map[string]string{    
			"message": "Method not allowed. Please use POST.",    
		})    
		return    
	}    
    
	var userInput model.UserInput    
    
	// Parse JSON request body    
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {    
		log.Printf("Invalid request data: %v", err)    
		w.Header().Set("Content-Type", "application/json")    
		w.WriteHeader(http.StatusBadRequest)    
		json.NewEncoder(w).Encode(map[string]string{    
			"message": "Invalid input. Please check your request data.",    
		})    
		return    
	}    
    
	// Find user by email    
	var user model.Users    
	if err := config.DB.Where("email = ?", userInput.Email).First(&user).Error; err != nil {    
		log.Printf("User not found: %v", err)    
		w.Header().Set("Content-Type", "application/json")    
		w.WriteHeader(http.StatusUnauthorized)    
		json.NewEncoder(w).Encode(map[string]string{    
			"message": "Invalid email or password.",    
		})    
		return    
	}    
    
	// Compare password    
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {    
		log.Printf("Invalid password: %v", err)    
		w.Header().Set("Content-Type", "application/json")    
		w.WriteHeader(http.StatusUnauthorized)    
		json.NewEncoder(w).Encode(map[string]string{    
			"message": "Invalid email or password.",    
		})    
		return    
	}    
    
	// Generate JWT token    
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{    
		"id":    user.ID,    
		"email": user.Email,    
		"role":  user.RoleID,    
		"exp":   time.Now().Add(24 * time.Hour).Unix(),    
	})    
    
	tokenString, err := token.SignedString([]byte(config.JwtKey)) // Pastikan config.JwtKey sudah diatur    
	if err != nil {    
		log.Printf("Failed to generate token: %v", err)    
		w.Header().Set("Content-Type", "application/json")    
		w.WriteHeader(http.StatusInternalServerError)    
		json.NewEncoder(w).Encode(map[string]string{    
			"message": "Failed to generate token.",    
		})    
		return    
	}    
    
	// Set expiration time for the token    
	expirationTime := time.Now().Add(24 * time.Hour)    
    
	// Create a new Tokens entry    
	newToken := model.Tokens{    
		UserID:    user.ID,    
		Token:     tokenString,    
		CreatedAt: time.Now(),    
		ExpiresAt: expirationTime,    
	}    
    
	// Save the token to the database    
	if err := config.DB.Create(&newToken).Error; err != nil {    
		log.Printf("Failed to save token to database: %v", err)    
		w.Header().Set("Content-Type", "application/json")    
		w.WriteHeader(http.StatusInternalServerError)    
		json.NewEncoder(w).Encode(map[string]string{    
			"message": "Failed to save token to database.",    
		})    
		return    
	}    
    
	// Create a new ActiveTokens entry    
	newActiveToken := model.ActiveTokens{    
		TokenID:   newToken.ID,    
		CreatedAt: time.Now(),    
	}    
    
	// Save the active token to the database    
	if err := config.DB.Create(&newActiveToken).Error; err != nil {    
		log.Printf("Failed to save active token to database: %v", err)    
		w.Header().Set("Content-Type", "application/json")    
		w.WriteHeader(http.StatusInternalServerError)    
		json.NewEncoder(w).Encode(map[string]string{    
			"message": "Failed to save active token to database.",    
		})    
		return    
	}    
    
	// Send success response with token    
	w.Header().Set("Content-Type", "application/json")    
	w.WriteHeader(http.StatusOK)    
	json.NewEncoder(w).Encode(map[string]interface{}{    
		"message": "Login successful",    
		"token":   tokenString,    
	})    
}  
