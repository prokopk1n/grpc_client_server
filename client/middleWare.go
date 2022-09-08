package main

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)

func checkAuthMiddleware(db *sql.DB, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, ok := r.Cookie("jwt-token")
		if ok != nil {
			if ok == http.ErrNoCookie {
				log.Printf("No cookie: %v", ok)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			log.Printf("Check auth: %v", ok)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tknStr := cookie.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}
