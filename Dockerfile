FROM golang:1.19.1-alpine3.15 as builder
COPY go.mod go.sum /go/src/client_server/
WORKDIR /go/src/client_server/client/src
RUN go mod download
COPY . /go/src/client_server
RUN go build -o ./build /go/src/client_server/client/src/main.go
EXPOSE 8080:8080
ENTRYPOINT ["./build"]