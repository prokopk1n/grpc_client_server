package main

import (
	"client_server/session"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"html/template"
	"log"
	"net/http"
	"time"
)

type TicketInfo struct {
	FlightDate    time.Time
	PassengerName string
	FlightFrom    string
	FlightTo      string
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
			log.Fatalf("Can not get ticket id = %s info: %v\n", ticketNumber, err)
		}
		ticketInfo := TicketInfo{PassengerName: reply.GetPassengerName(), FlightFrom: reply.GetFlightFrom(),
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

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	w.Header().Set("Accept-Charset", "utf-8")

	tmpl := template.Must(template.ParseFiles("client/html/mainpage.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Can not execute template: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ticketPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/html/ticket.html")
}

func main() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("can not connect to localhost")
	}
	defer conn.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainPageHandler)
	mux.HandleFunc("/ticketinfo", ticketPageHandler)
	mux.HandleFunc("/ticket", handlerSlashCreate(conn))
	mux.Handle("/styles.css", http.FileServer(http.Dir("client/html")))
	fmt.Println("Client: start working...")
	err = http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatalln("can not listen port 8080: %v", err)
	}
}
