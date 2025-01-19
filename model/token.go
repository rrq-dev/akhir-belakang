package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Tokens model
type Tokens struct {  
	ID        int       `gorm:"primaryKey" json:"id"`  
	UserID    int       `json:"user_id"`  
	Token     string    `json:"token"`  
	CreatedAt time.Time `json:"created_at"`  
	ExpiresAt time.Time `json:"expires_at"`  
}  
  
// ActiveTokens model  
type ActiveTokens struct {  
	ID        int       `gorm:"primaryKey" json:"id"`  
	TokenID   int       `json:"token_id"`  
	CreatedAt time.Time `json:"created_at"`  
}  
  
// BlacklistTokens model  
type BlacklistTokens struct {    
	ID        int        `gorm:"primaryKey" json:"id"`    
	Token     int        `json:"unique;not null"` // Store the ID of the blacklisted token    
	ExpiresAt time.Time  `json:"expires_at"` // Expiration time of the token    
	CreatedAt time.Time  `json:"created_at"` // Creation time of the blacklist entry    
}  

  
// Claims model   
type Claims struct {
    UserID    int       `json:"user_id"`
    Role      string    `json:"role"`
    ExpiresAt time.Time `json:"expires_at"`
    jwt.StandardClaims
}

func (c *Claims) UnmarshalJSON(data []byte) error {
    type Alias Claims
    aux := &struct {
        Role interface{} `json:"role"`
        *Alias
    }{
        Alias: (*Alias)(c),
    }

    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }

    switch v := aux.Role.(type) {
    case string:
        c.Role = v
    case float64:
        c.Role = fmt.Sprintf("%v", v)
    default:
        return fmt.Errorf("invalid type for role: %T", v)
    }

    return nil
}


