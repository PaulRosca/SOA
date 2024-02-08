package function

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Type  string `json:"type"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("my_secret_key")

func Handle(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Cannot parse token", http.StatusBadRequest)
		return
	}
	tokenString := strings.Split(tokenCookie.String(), "token=")[1]
	claim := Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claim, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !token.Valid {
		http.Error(w, "Invalid Tokan!", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(claim)
}
