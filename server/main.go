package main

import (
	"client_server/session"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"
)

const (
	postgresConn = "user=postgres password=cskamoscoW99 dbname=demo sslmode=disable"
)

type dataBase struct {
	*sql.DB
}

func (db *dataBase) GetTicketId(id string) (*session.TicketInfoReply, error) {
	rows, err := db.Query("SELECT passenger_name, flight_id, scheduled_departure, departure_airport, arrival_airport "+
		"FROM (Tickets JOIN Ticket_flights USING (ticket_no)) JOIN Flights USING (flight_id) WHERE ticket_no = $1;", id)
	if err != nil {
		log.Printf("Error while query to database: %v", err)
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		log.Printf("No result for %s", id)
		return nil, fmt.Errorf("Not found")
	}
	reply := session.TicketInfoReply{}
	var flightDate time.Time
	err = rows.Scan(&reply.PassengerName, &reply.FlightId, &flightDate, &reply.FlightFrom, &reply.FlightTo)
	reply.FlightDate = timestamppb.New(flightDate)
	if err != nil {
		log.Printf("can not parse result to TicketInfoReply: %v", err)
		return nil, err
	}
	return &reply, nil
}

type server struct {
	db *dataBase
	session.UnimplementedAirplaneServerServer
}

func (s *server) GetTicketInfo(ctx context.Context, in *session.TicketReq) (*session.TicketInfoReply, error) {
	flightId := in.GetTicketNo()
	res, err := s.db.GetTicketId(flightId)
	if err != nil && err.Error() == "Not found" {
		return nil, nil
	}
	return res, err
}

func main() {
	db, err := sql.Open("postgres", postgresConn)
	if err != nil {
		log.Fatalf("Can not connect to database: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalln("Can not listen port 8081")
	}
	newServer := grpc.NewServer()
	session.RegisterAirplaneServerServer(newServer, &server{db: &dataBase{db}})
	log.Printf("Server listening port 8081")
	if err := newServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
