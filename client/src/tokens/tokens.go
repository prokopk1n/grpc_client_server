package tokens

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
)

var JwtKey = []byte("secret_key")

type Claims struct {
	UserId int
	jwt.StandardClaims
}

func GetUserId(jwtToken string) (int, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(jwtToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err == nil {
		return claims.UserId, nil
	}
	if v, ok := err.(*jwt.ValidationError); ok {
		if v.Errors == jwt.ValidationErrorExpired {
			return claims.UserId, nil
		}
	}
	return 0, err
}

func CheckRefreshToken(r *http.Request, db *sql.DB) error {
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

func NewJWTToken(userId int, expireTime time.Time) (string, error) {
	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Printf("NEW JWT Token = %+v", token.Claims)
	return token.SignedString(JwtKey)
}

func NewRefreshToken(expireTime time.Time, db *sql.DB) (string, error) {
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

func NewTokenPair(jwtLifeTime, refreshLifeTime time.Duration, userId int, db *sql.DB) (jwtToken string, refreshToken string, err error) {
	jwtToken, err = NewJWTToken(userId, time.Now().Add(jwtLifeTime))
	if err != nil {
		log.Printf("Can not create JWT token: %v", err)
		return
	}
	refreshToken, err = NewRefreshToken(time.Now().Add(refreshLifeTime), db)
	if err != nil {
		log.Printf("Can not create Refresh token: %+v", err)
		return
	}
	return
}
