package middleware

import (
	"encoding/json"
	"myrest-api/pkg/model"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type AuthenticationMiddleware struct {
}

func (amw *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/register" || r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
		} else {

			cookie, _ := r.Cookie("token")
			if cookie.Value == "some_token_value" {
				next.ServeHTTP(w, r)
			} else {
				var answer model.Answer
				answer.Message = "Wrong token"
				json.NewEncoder(w).Encode(answer)
			}
		}
	})
}
