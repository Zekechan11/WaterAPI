package middleware

import (
    "net/http"
    "strings"
    "context"
	"water-api/util"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        claims, err := util.ValidateToken(tokenString)
        if err != nil {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        
        ctx := context.WithValue(r.Context(), "username", claims.Username)
        r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}
