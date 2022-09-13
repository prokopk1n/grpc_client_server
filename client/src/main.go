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
	postgresConn = "host=host.docker.internal port=5431 user=postgres password=postgres dbname=clientdb sslmode=disable"
)

func main() {
	conn, err := grpc.Dial("host.docker.internal:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("can not connect to localhost")
	}
	defer conn.Close()

	db, err := sql.Open("postgres", postgresConn)
	if err != nil {
		log.Fatalf("Can not connect to : %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Can not ping db: %v", err)
	}
	defer db.Close()

	handlerManager := handlers.NewHandlerManager(conn, db, time.Second*10, time.Hour)
	mux := handlerManager.Init()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}
	fmt.Println("Client: start working...")
	err = srv.ListenAndServeTLS("../key/server.crt", "../key/server.key")
	if err != nil {
		log.Fatalln("can not listen port 8080: %v", err)
	}
}
