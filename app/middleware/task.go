package middleware

import (
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var cookieToken string // JWT-токен из куки
			cookie, err := r.Cookie("token")
			if err == nil {
				cookieToken = cookie.Value
			}
			jwtInstance := jwt.New(jwt.SigningMethodHS256)
			token, err := jwtInstance.SignedString([]byte(pass))

			if cookieToken != token {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
