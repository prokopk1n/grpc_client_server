package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"time"
)

const (
	postgresConn = "user=postgres password=password dbname=clientdb sslmode=disable"
)

var db *sql.DB

type TicketInfo struct {
	FlightDate    time.Time
	PassengerName string
	FlightFrom    string
	FlightTo      string
}

func main() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("can not connect to localhost")
	}
	defer conn.Close()

	db, err = sql.Open("postgres", postgresConn)
	if err != nil {
		log.Fatalf("Can not connect to : %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", MainPageHandler)
	mux.HandleFunc("/ticketinfo", TicketPageHandler)
	mux.HandleFunc("/ticket", HandlerSlashCreate(conn))
	mux.HandleFunc("/signin", HandlerLogin)
	mux.HandleFunc("/lk", CheckAuthMiddleware(LKhandler))
	mux.HandleFunc("/check_auth", CheckAuth)
	mux.HandleFunc("/signup", HandlerSignUp)
	mux.HandleFunc("/addticket", CheckAuthMiddleware(HandlerAddTicket))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./template/css"))))
	fmt.Println("Client: start working...")
	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}
	err = srv.ListenAndServeTLS("key/server.crt", "key/server.key")
	if err != nil {
		log.Fatalln("can not listen port 8080: %v", err)
	}
}
