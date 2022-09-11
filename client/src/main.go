package main

import (
	"client_server/client/src/handlers"
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
var conn *grpc.ClientConn

func main() {
	var err error
	conn, err = grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("can not connect to localhost")
	}
	defer conn.Close()

	db, err = sql.Open("postgres", postgresConn)
	if err != nil {
		log.Fatalf("Can not connect to : %v", err)
	}
	defer db.Close()

	handlerManager := handlers.NewHandlerManager(conn, db, time.Second*10, time.Hour)
	mux := handlerManager.Init()

	fmt.Println("Client: start working...")
	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}
	err = srv.ListenAndServeTLS("../key/server.crt", "../key/server.key")
	if err != nil {
		log.Fatalln("can not listen port 8080: %v", err)
	}
}
