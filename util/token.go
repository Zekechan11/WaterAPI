package util

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your_secret_key")

type Claims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func GenerateToken(username string) (string, error) {
    claims := Claims{
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    
    if err != nil || !token.Valid {
        return nil, err
    }
    return claims, nil
}
