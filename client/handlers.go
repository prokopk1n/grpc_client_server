package main

import (
	"client_server/session"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"html/template"
	"log"
	"net/http"
	"time"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func verifyUserPass(user string, pass string, db *sql.DB) bool {
	rows, err := db.Query("SELECT passw FROM users WHERE email = $1", user)
	if err != nil {
		log.Printf("Can not find user %s in DB: %v", user, err)
		return false
	}
	if !rows.Next() {
		log.Printf("Empty result with login %s", user)
		return false
	}
	var password string
	err = rows.Scan(&password)
	if err != nil {
		log.Printf("Can not scan password: %v", err)
		return false
	}
	if compared := bcrypt.CompareHashAndPassword([]byte(password), []byte(pass)); compared == nil {
		log.Printf("Auth was passed successfully")
		return true
	} else {
		log.Printf("Auth was not passed: %v", compared)
		return false
	}
}

func ticketPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/ticket.html")
}

func LKhandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("YOU ARE WELCOME"))
}

func checkAuth(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Can not parse form: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	email := r.FormValue("auth_email")
	password := r.FormValue("auth_pass")

	if verifyUserPass(email, password, db) {
		claims := &Claims{
			Email:          email,
			StandardClaims: jwt.StandardClaims{},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			log.Printf("Can not create token: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cookie := http.Cookie{
			Name:  "jwt-token",
			Value: tokenString,
			Path:  "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/lk", http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Неверный логин или пароль"))
	}

}

func handlerSlashCreate(conn *grpc.ClientConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get request")
		client := session.NewAirplaneServerClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		params := r.URL.Query()
		if !params.Has("id") {
			w.WriteHeader(404)
			return
		}
		ticketNumber := params["id"][0]
		reply, err := client.GetTicketInfo(ctx, &session.TicketReq{TicketNo: ticketNumber})
		if err != nil {
			log.Printf("Can not get ticket id = %s info: %v\n", ticketNumber, err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Can not get ticket id = %s", ticketNumber)
			return
		}
		ticketInfo := TicketInfo{FlightDate: reply.GetFlightDate().AsTime(), PassengerName: reply.GetPassengerName(), FlightFrom: reply.GetFlightFrom(),
			FlightTo: reply.GetFlightTo()}
		//ticketInfo.FlightDate = time.
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(ticketInfo)
		if err != nil {
			log.Fatalln("Can not encode to json format")
		}
	}
}

func handlerSignIn(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/auth.html")
}

func handlerSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "html/registration.html")
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Printf("Can not parse form Sign Up: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	email := r.FormValue("email")
	password, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 0)
	if err != nil {
		log.Printf("Can not generate hash frow password: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = db.Query("INSERT INTO users(email, passw) VALUES ($1, $2)", email, password)
	if err != nil {
		log.Printf("Can not create new user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/lk", http.StatusSeeOther)
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("main page handler")
	w.Header().Set("Content-type", "text/html")
	w.Header().Set("Accept-Charset", "utf-8")

	tmpl := template.Must(template.ParseFiles("html/mainpage.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Can not execute template: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
