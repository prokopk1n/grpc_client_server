package server

import (
	"client_server/session"
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

type DataBase struct {
	*sql.DB
}

func NewDataBase(db *sql.DB) *DataBase {
	return &DataBase{db}
}

func (db *DataBase) GetTicketId(id string) (*session.TicketInfoReply, error) {
	rows, err := db.Query("SELECT passenger_name, flight_id, scheduled_departure, departure_airport, arrival_airport "+
		"FROM (Tickets JOIN Ticket_flights USING (ticket_no)) JOIN Flights USING (flight_id) WHERE ticket_no = $1", id)
	if err != nil {
		log.Printf("Error while query to database: %v\n", err)
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

type Server struct {
	db *DataBase
	session.UnimplementedAirplaneServerServer
}

func NewServer(db *DataBase) *Server {
	return &Server{db: db}
}

func (s *Server) GetTicketInfo(ctx context.Context, in *session.TicketReq) (*session.TicketInfoReply, error) {
	flightId := in.GetTicketNo()
	res, err := s.db.GetTicketId(flightId)
	if err != nil && err.Error() == "Not found" {
		return nil, nil
	}
	return res, err
}
