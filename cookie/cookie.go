package cookie

import (
	"net/http"
	"time"
)

func Add(w http.ResponseWriter, name, value string) {
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: time.Now().AddDate(3, 0, 0),
	}

	http.SetCookie(w, &cookie)
}

func HasUser(r *http.Request) bool {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "sessionID" {
			return true
		}
	}

	return false
}

func GetUserID(r *http.Request) string {
	sessionCookie, err := r.Cookie("sessionID")
	if err != nil {
		panic(err.Error())
	}

	return sessionCookie.Value
}
