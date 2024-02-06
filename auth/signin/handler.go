package function

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

var jwtSecret = []byte("my_secret_key")

func init() {
	var err error
	cfg := mysql.Config{
		User:   "root",
		Passwd: "test1234",
		Net:    "tcp",
		Addr:   "mysql.default.svc.cluster.local:3306",
		DBName: "auth",
	}
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type User struct {
	Credentials
	ID         int64  `json:"id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}

type Claims struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func Handle(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var user User
	row := db.QueryRow("SELECT * FROM user WHERE email = ?", credentials.Email)
	if err := row.Scan(&user.ID, &user.Email, &user.First_name, &user.Last_name, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password!", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid email or password!", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		ID:    user.ID,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = ""
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
