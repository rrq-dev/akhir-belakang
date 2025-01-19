package model

import "time"

// Lokasi model
type Lokasi struct {  
    ID          int       `gorm:"primaryKey" json:"id"`  
    Nama        string    `json:"nama"`  
    Alamat      string    `json:"alamat"`  
    Deskripsi   string    `json:"deskripsi"`  
    CreatedAt   time.Time `json:"created_at"`  
} 
// Override default table name
func (Lokasi) TableName() string {
    return "lokasi"
} 
  
// Role model  
type Role struct {  
    ID        int       `gorm:"primaryKey" json:"id"`  
    Nama      string    `json:"nama"`  
    CreatedAt time.Time `json:"created_at"`  
}  
  
// Users model  
type Users struct {  
	ID        int       `gorm:"primaryKey" json:"id"`  
	Username  string    `json:"username"`  
	Password  string    `json:"password"`  
	Email     string    `json:"email"`  
	RoleID    int       `json:"role_id"`  
	CreatedAt time.Time `json:"created_at"`  
}  
  
// Feedback model  
type Feedback struct {  
    ID        int       `gorm:"primaryKey" json:"id"`  
    UserID    int       `json:"user_id"`  
    LokasiID  int       `json:"lokasi_id"`  
    Komentar  string    `json:"komentar"`  
    Rating    int       `json:"rating"`  
    CreatedAt time.Time `json:"created_at"`  
}  
  
// Tokens model  
type UserInput struct {  
	Email           string `json:"email" binding:"required,email"`  
	Username        string `json:"username" binding:"required"`  
	Password        string `json:"password" binding:"required"`  
	ConfirmPassword string `json:"confirm_password" binding:"required"`  
}

// Di dalam model  
type ContextKey string  
  
const (  
	UserIDKey ContextKey = "userID"  
	RoleKey   ContextKey = "role"  
)  
