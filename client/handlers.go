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

func verifyUserPass(user string, pass string, db *sql.DB) (int, bool) {
	rows, err := db.Query("SELECT user_id, passw FROM users WHERE email = $1", user)
	if err != nil {
		log.Printf("Can not find user %s in DB: %v", user, err)
		return 0, false
	}
	if !rows.Next() {
		log.Printf("Empty result with login %s", user)
		return 0, false
	}
	var password string
	var userId int
	err = rows.Scan(&userId, &password)
	if err != nil {
		log.Printf("Can not scan password: %v", err)
		return 0, false
	}
	if compared := bcrypt.CompareHashAndPassword([]byte(password), []byte(pass)); compared == nil {
		log.Printf("Auth was passed successfully")
		return userId, true
	} else {
		log.Printf("Auth was not passed: %v", compared)
		return 0, false
	}
}

func TicketPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/ticket.html")
}

func findAllTickets(jwtToken string) ([]*TicketInfo, error) {
	userId, err := getUserId(jwtToken)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT ticket FROM userstickets WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	data := make([]*TicketInfo, 0, 10)
	var ticketId string
	for rows.Next() {
		err = rows.Scan(&ticketId)
		if err != nil {
			return nil, err
		}
		client := session.NewAirplaneServerClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		req := &session.TicketReq{TicketNo: ticketId}
		if client == nil {
			log.Printf("CLIENT IS NIL")
		}
		reply, err := client.GetTicketInfo(ctx, req)
		if err != nil {
			data = append(data, &TicketInfo{TicketId: ticketId, FlightDate: time.Time{}, PassengerName: "error", FlightFrom: "error",
				FlightTo: "error"})
		} else {
			data = append(data, &TicketInfo{TicketId: ticketId, FlightDate: reply.GetFlightDate().AsTime(), PassengerName: reply.GetPassengerName(), FlightFrom: reply.GetFlightFrom(),
				FlightTo: reply.GetFlightTo()})
		}
	}
	return data, nil
}

func LKhandler(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("template/lk.html", "template/head.html")
	if err != nil {
		log.Printf("Can not parse lk template: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	value, err := r.Cookie("jwt-token")
	data, err := findAllTickets(value.Value)
	if err != nil {
		log.Printf("Error while get all tickets from server: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(data) == 0 {
		data = nil
	}
	err = templ.ExecuteTemplate(w, "lk", data)
	if err != nil {
		log.Printf("Can not execute template lk with data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Can not parse form: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	email := r.FormValue("auth_email")
	password := r.FormValue("auth_pass")

	if userId, ok := verifyUserPass(email, password, db); ok {
		jwtToken, refreshToken, err := newTokenPair(userId)
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
		http.Redirect(w, r, "/lk", http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Неверный логин или пароль"))
	}

}

func HandlerSlashCreate(conn *grpc.ClientConn) func(http.ResponseWriter, *http.Request) {
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
		ticketInfo := TicketInfo{TicketId: ticketNumber, FlightDate: reply.GetFlightDate().AsTime(), PassengerName: reply.GetPassengerName(), FlightFrom: reply.GetFlightFrom(),
			FlightTo: reply.GetFlightTo()}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(ticketInfo)
		if err != nil {
			log.Fatalln("Can not encode to json format")
		}
	}
}

func HandlerAddTicket(w http.ResponseWriter, r *http.Request) {
	var claims Claims
	cookie, _ := r.Cookie("jwt-token")
	jwt.ParseWithClaims(cookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	_, err := db.Query("INSERT INTO userstickets(user_id, ticket) VALUES($1, $2)", claims.UserId, r.FormValue("ticket_id"))
	if err != nil {
		log.Printf("Can not create new ticket for user %d: %v", claims.UserId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/lk", http.StatusSeeOther)
}

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/login.html")
}

func HandlerSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "template/registration.html")
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
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("main page handler")
	w.Header().Set("Content-type", "text/html")
	w.Header().Set("Accept-Charset", "utf-8")

	tmpl, err := template.ParseFiles("template/main_page.html", "template/main.html", "template/head.html")
	if err != nil {
		log.Printf("Can not parse main page template: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "main_page", nil)
	if err != nil {
		log.Printf("Can not execute template: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
