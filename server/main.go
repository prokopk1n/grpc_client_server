package main

import (
	"client_server/server/server"
	"client_server/session"
	"database/sql"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	postgresConn = "user=postgres password=password dbname=demo sslmode=disable"
)

func main() {
	db, err := sql.Open("postgres", postgresConn)
	if err != nil {
		log.Fatalf("Can not connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Query("SET lc_messages TO 'en_US.UTF-8'")
	if err != nil {
		log.Printf("Can not change locale of postgres to English")
	}

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("Can not listen port 8081")
	}
	newServer := grpc.NewServer()
	session.RegisterAirplaneServerServer(newServer, server.NewServer(server.NewDataBase(db)))
	log.Printf("Server listening port 8081")
	if err := newServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
