package middleware

import (
	"log"
	"net/http"
)

type AuthenticationMiddleware struct {
	TokenUsers map[string]string
}

func (amw *AuthenticationMiddleware) Populate() {
	amw.TokenUsers["00000000"] = "admin"
	amw.TokenUsers["aaaaaaaa"] = "userA"
	amw.TokenUsers["05f717e5"] = "randomUser"
	amw.TokenUsers["deadbeef"] = "user0"
}

func (amw *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")

		if user, found := amw.TokenUsers[token]; found {
			// Мы нашли токен в нашей карте
			log.Printf("Authenticated user %s\n", user)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
