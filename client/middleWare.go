package main

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)

func CheckAuthMiddleware(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, ok := r.Cookie("jwt-token")
		if ok != nil {
			if ok == http.ErrNoCookie {
				log.Printf("No cookie: %v", ok)
			}
			log.Printf("Check auth: %v", ok)
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}
		var claims Claims
		_, err := jwt.ParseWithClaims(cookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			} else if v, ok := err.(*jwt.ValidationError); ok {
				if v.Errors == jwt.ValidationErrorExpired {
					log.Printf("Token jwt lifetime expired")
					err = checkRefreshToken(r)
					if err != nil {
						http.Redirect(w, r, "/signin", http.StatusSeeOther)
						return
					}
					jwtToken, refreshToken, err := newTokenPair(claims.UserId)
					if err != nil {
						log.Printf("Error occurs in newTokenPair(): %v", err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					log.Printf("Created new pair of tokens")
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
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}
		}
		log.Printf("Auth was passed successfully")
		handler(w, r)
	}
}
