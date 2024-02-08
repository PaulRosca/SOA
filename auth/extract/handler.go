package function

import (
	"encoding/json"
	"io"
	"net/http"

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
	byteData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokenString := string(byteData)
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
