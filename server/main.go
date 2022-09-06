package main

import (
	"client_server/session"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type dataBase struct {
	db map[string]map[string]string
}

func (d *dataBase) GetTicketInfo(id string) (*session.TicketInfoReply, error) {
	if _, ok := d.db[id]; !ok {
		return nil, fmt.Errorf("can not find id = %s", id)
	}
	reply := session.TicketInfoReply{
		FlightDate:    "",
		FlightId:      id,
		PassengerName: d.db[id]["name"],
		FlightFrom:    d.db[id]["from"],
		FlightTo:      d.db[id]["to"],
	}
	return &reply, nil
}

type server struct {
	session.UnimplementedAirplaneServerServer
}

var db = dataBase{db: map[string]map[string]string{"123": map[string]string{"name": "Sergey", "from": "London", "to": "Moscow"}}}

func (s *server) GetTicketInfo(ctx context.Context, in *session.TicketReq) (*session.TicketInfoReply, error) {
	flightId := in.GetTicketNo()
	return db.GetTicketInfo(flightId)
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalln("Can not listen port 8081")
	}
	newServer := grpc.NewServer()
	session.RegisterAirplaneServerServer(newServer, &server{})
	log.Printf("Server listening port 8081")
	if err := newServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
