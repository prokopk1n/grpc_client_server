package main

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)

func CheckAuthMiddleware(db *sql.DB, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
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
		_, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if v, ok := err.(*jwt.ValidationError); ok {
				if v.Errors == jwt.ValidationErrorExpired {
					err = checkRefreshToken(r)
					if err != nil {
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
					jwtToken, refreshToken, err := newTokenPair()
					if err != nil {
						log.Printf("Error occurs in newTokenPair(): %v", err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					http.SetCookie(w, &http.Cookie{
						Name:     "jwt-token",
						Value:    jwtToken,
						Path:     "/",
						HttpOnly: true,
					})
					http.SetCookie(w, &http.Cookie{
						Name:     "refresh-token",
						Value:    refreshToken,
						Path:     "/",
						HttpOnly: true,
					})
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		log.Printf("Auth was passed successfully")
		handler(w, r)
	}
}
