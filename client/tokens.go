package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	jwt.StandardClaims
}

func checkRefreshToken(r *http.Request) error {
	cookie, ok := r.Cookie("refresh-token")
	if ok != nil {
		log.Printf("Can not find refresh-token in cookie")
		return ok
	}
	result, err := db.Query("SELECT expire_time FROM refresh_tokens WHERE token = $1", []byte(cookie.Value))
	if err != nil {
		log.Printf("Error occurs when query to db to find refresh token: %v", err)
		return err
	}
	defer result.Close()

	if !result.Next() {
		return fmt.Errorf("Can not find refresh token")
	}

	db.Query("DELETE FROM refresh_tokens WHERE token = $1", []byte(cookie.Value))

	var expireTime int64
	err = result.Scan(&expireTime)
	if err != nil {
		log.Printf("Can not scan result: %v", err)
		return err
	}

	if time.Unix(expireTime, 0).Before(time.Now()) {
		return fmt.Errorf("refresh token was expired")
	}
	return nil
}

func newJWTToken(expireTime time.Time) (string, error) {
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func newRefreshToken(expireTime time.Time) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	refreshToken := base64.URLEncoding.EncodeToString(b)
	_, err := db.Query("INSERT INTO refresh_tokens(token, expire_time) VALUES($1, $2)", []byte(refreshToken), expireTime.Unix())
	if err != nil {
		log.Printf("Can not insert new refresh token into db: %v", err)
		return "", err
	}
	return refreshToken, nil
}

func newTokenPair() (jwtToken string, refreshToken string, err error) {
	jwtToken, err = newJWTToken(time.Now().Add(time.Second * 10))
	if err != nil {
		log.Printf("Can not create JWT token: %v", err)
		return
	}
	refreshToken, err = newRefreshToken(time.Now().Add(time.Hour * 72))
	if err != nil {
		log.Printf("Can not create Refresh token: %+v", err)
		return
	}
	return
}
