package main

import (
	"client_server/client/src/handlers"
	"client_server/server/server"
	"client_server/session"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	dbUsers      *sql.DB
	dbTickets    *sql.DB
	err          error
	client       = &http.Client{Timeout: time.Second * 5}
	clientServer *httptest.Server
)

const (
	dbTicketsConn = "user=postgres password=postgres port=5431 dbname=demo sslmode=disable"
	dbUsersConn   = "user=postgres password=postgres port=5431 dbname=clientdb sslmode=disable"
)

func TestInit(t *testing.T) {
	dbTickets, err = sql.Open("postgres", dbTicketsConn)
	if err != nil {
		t.Fatalf("Error when open db occurs: %v", err)
	}

	err = dbTickets.Ping()
	if err != nil {
		t.Fatalf("Can not ping db: %v", err)
	}

	newServer := grpc.NewServer()
	session.RegisterAirplaneServerServer(newServer, server.NewServer(server.NewDataBase(dbTickets)))
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		t.Fatalf("Can not listen port 8081: %v", err)
	}

	go func() {
		if err := newServer.Serve(lis); err != nil {
			t.Fatalf("Error grpc: %v", err)
		}
	}()

	dbUsers, err = sql.Open("postgres", dbUsersConn)
	if err != nil {
		t.Fatalf("Can not open db users: %v", err)
	}

	err = dbUsers.Ping()
	if err != nil {
		t.Fatalf("Can not ping db users: %v", err)
	}

	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("can not connect to localhost")
	}
	handler := handlers.NewHandlerManager(conn, dbUsers, time.Second*10, time.Hour)
	clientServer = httptest.NewServer(handler.Init())
	runTests(t)
}

func runTests(t *testing.T) {
	TestFindTicketById(t)
}

func TestFindTicketById(t *testing.T) {
	cases := []struct {
		path   string
		query  string
		status int
		result *handlers.TicketInfo
	}{
		{
			path:   "/ticket",
			query:  "id=0005432000987",
			status: http.StatusOK,
			result: &handlers.TicketInfo{
				TicketId:      "0005432000987",
				FlightDate:    time.Date(2017, 7, 16, 9, 5, 0, 0, time.UTC),
				PassengerName: "VALERIY TIKHONOV",
				FlightFrom:    "CSY",
				FlightTo:      "SVO",
			},
		},
		{
			path:   "/ticket",
			query:  "id=1115002000987",
			status: http.StatusBadRequest,
			result: nil,
		},
	}
	for _, value := range cases {
		caseName := fmt.Sprintf("Case: %+v", value)
		if clientServer == nil {
			fmt.Printf("Client server is nil")
		}
		req, _ := http.NewRequest("GET", clientServer.URL+value.path+"?"+value.query, nil)
		//fmt.Printf("New request: %s\n", req.URL)
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("[%s] request error: %v", caseName, err)
			continue
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != value.status {
			t.Errorf("expected http status: %d got: %d", value.status, resp.StatusCode)
			continue
		}
		if value.status != http.StatusBadRequest {
			var result handlers.TicketInfo
			err = json.Unmarshal(body, &result)
			if err != nil {
				t.Errorf("cant unpack json: %v", err)
				continue
			}
			if result.TicketId != value.result.TicketId ||
				result.FlightDate != value.result.FlightDate ||
				result.FlightTo != value.result.FlightTo ||
				result.FlightFrom != value.result.FlightFrom ||
				result.PassengerName != value.result.PassengerName {
				t.Errorf("results not match\nGot: %#v\nExpected: %#v", result, value.result)
				continue
			}
		}
	}
}
