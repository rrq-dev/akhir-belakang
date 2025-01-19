package helper

import (
	"akhir-belakang/config"
	"akhir-belakang/model"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/golang-jwt/jwt"
)

// GenerateToken generates a new JWT token
func GenerateToken(userID int, role string) (string, error) {  
	claims := &model.Claims{  
		UserID:    userID,  
		Role:      role, // Ensure this is a string  
		ExpiresAt: time.Now().Add(time.Hour * 24).UTC(), // Set expiration time  
		StandardClaims: jwt.StandardClaims{  
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Set expiration time in Unix format  
		},  
	}  
  
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)  
	return token.SignedString([]byte(config.JwtKey)) // Sign the token with your secret key  
}     
  
func BlacklistToken(w http.ResponseWriter, r *http.Request) {  
	tokenString, err := GetTokenFromHeader(r)  
	if err != nil {  
		log.Printf("Token error: %v", err)  
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)  
		return  
	}  
  
	// Verifikasi token JWT  
	claims := &model.Claims{}  
	if err := ParseAndValidateToken(tokenString, claims); err != nil {  
		log.Printf("Token validation error: %v", err)  
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)  
		return  
	}  
  
	// Simpan token ke dalam database sebagai blacklist  
	blacklistToken := model.BlacklistTokens{  
		Token:     claims.UserID, // Use UserID as the token ID for blacklisting  
		ExpiresAt: time.Now().Add(24 * time.Hour), // Set expiration time for the blacklist entry  
		CreatedAt: time.Now(),  
	}  
  
	if err := config.DB.Create(&blacklistToken).Error; err != nil {  
		log.Printf("Failed to insert token into blacklist: %v", err)  
		http.Error(w, "Failed to blacklist token", http.StatusInternalServerError)  
		return  
	}  
  
	// Kirim respons logout berhasil  
	WriteResponse(w, http.StatusOK, map[string]string{  
		"message": "Logout successful, token has been blacklisted",  
	})  
}  
  
// ValidateTokenMiddleware validates the JWT token and inserts user info into the context
func ValidateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from header
		tokenString, err := GetTokenFromHeader(r)
		if err != nil {
			log.Printf("Token retrieval error: %v", err)
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Check if the token is blacklisted
		if IsTokenBlacklisted(tokenString) {
			log.Printf("Token is blacklisted: %v", tokenString)
			http.Error(w, "Unauthorized: Token has been blacklisted", http.StatusUnauthorized)
			return
		}

		// Parse and validate the JWT token
		claims := &model.Claims{}
		if err := ParseAndValidateToken(tokenString, claims); err != nil {
			log.Printf("Token validation failed: %v", err)
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Log successful validation
		log.Printf("Token validated successfully: userID=%d, role=%s", claims.UserID, claims.Role)

		// Insert userID and role into request context
		ctx := context.WithValue(r.Context(), model.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, model.RoleKey, claims.Role)

		// Pass the updated context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
 
  
// getTokenFromHeader retrieves the token from the Authorization header  
func GetTokenFromHeader(r *http.Request) (string, error) {  
	authHeader := r.Header.Get("Authorization")  
	if authHeader == "" {  
		return "", fmt.Errorf("missing Authorization header")  
	}  
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {  
		return authHeader[7:], nil  
	}  
	return "", fmt.Errorf("invalid Authorization header format")  
}  
  
// isTokenBlacklisted checks if the token is in the blacklist  
func IsTokenBlacklisted(token string) bool {  
	var blacklistedToken model.BlacklistTokens  
	err := config.DB.Where("token = ?", token).First(&blacklistedToken).Error  
	return err == nil // If no error, the token is blacklisted  
}  
  
// parseAndValidateToken parses and validates the JWT token  
func ParseAndValidateToken(tokenString string, claims *model.Claims) error {  
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {  
		return []byte(config.JwtKey), nil  
	})  
	return err  
} 
  
// MoveExpiredTokensToBlacklist moves expired tokens to the blacklist  
func MoveExpiredTokensToBlacklist() {  
	var expiredTokens []model.ActiveTokens  
	if err := config.DB.Where("created_at < ?", time.Now().Add(-24*time.Hour)).Find(&expiredTokens).Error; err != nil {  
		log.Printf("Failed to find expired tokens: %v", err)  
		return  
	}  
  
	// Move tokens to blacklist_tokens  
	for _, token := range expiredTokens {  
		blacklistToken := model.BlacklistTokens{  
			Token:     token.TokenID, // Use TokenID as an int  
			ExpiresAt: time.Now().Add(24 * time.Hour), // Set expiration time  
			CreatedAt: time.Now(),  
		}  
		if err := config.DB.Create(&blacklistToken).Error; err != nil {  
			log.Printf("Failed to move token to blacklist: %v", err)  
			continue  
		}  
  
		// Delete token from active_tokens  
		if err := config.DB.Delete(&token).Error; err != nil {  
			log.Printf("Failed to delete token from active_tokens: %v", err)  
		}  
	}  
  
	log.Printf("Moved %d expired tokens to blacklist.", len(expiredTokens))  
}  
  
// ScheduleTokenCleanup schedules the token cleanup  
func ScheduleTokenCleanup() error {  
	scheduler := gocron.NewScheduler(time.UTC)  
  
	// Jadwalkan cleanup setiap jam  
	_, err := scheduler.Every(1).Hour().Do(func() {  
		log.Println("Running token cleanup...")  
		MoveExpiredTokensToBlacklist()  
	})  
	if err != nil {  
		return err  
	}  
  
	// Mulai scheduler  
	scheduler.StartAsync()  
	return nil  
}  
