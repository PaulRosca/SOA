package function

import (
	"net/http"
	"time"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Unix(0, 0),
		Path:    "/",
	})
}
