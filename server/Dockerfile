FROM golang:1.19.1-alpine3.15 as builder
COPY go.mod go.sum /go/src/grpc_server/
WORKDIR /go/src/grpc_server/server
RUN go mod download
COPY . /go/src/grpc_server
RUN go build -o ./build /go/src/grpc_server/server/main.go
EXPOSE 8081
ENTRYPOINT ["./build"]