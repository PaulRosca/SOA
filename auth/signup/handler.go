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
var adminSecret = "root"

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

type User struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Password   string `json:"password,omitempty"`
	Secret     string `json:"secret,omtiemtpy"`
}

type Claims struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func Handle(w http.ResponseWriter, r *http.Request) {
	var user User
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	passwordHash, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Password = passwordHash
	result, err := db.Exec("INSERT INTO user (email, first_name, last_name, password) VALUES (?, ?, ?, ?)", user.Email, user.First_name, user.Last_name, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Secret == adminSecret {
		db.Exec("UPDATE user SET type = \"admin\" WHERE id = ?", id)
	}
	user.ID = id

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		ID:    id,
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
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	})
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
