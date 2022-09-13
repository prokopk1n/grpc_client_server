# Клиент-серверное приложение с использованием фреймворка gRPC.
Приложение предусмотрено для хранения авиабилетов или поиска информации по ним. 

## Table of Contents
* [General info](#general-info)
* [Technologies](#technologies)
* [Setup](#setup)

## General info
За основу взята тестовая база данных https://postgrespro.ru/education/demodb

## Setup
### Install and launch:
    git clone https://github.com/prokopk1n/grpc_client_server
    go run server/main.go
    cd client/src
    go run main.go
### Launch tests:
    cd docker
    docker-compose -p newtest -f docker-compose-test.yml -d up
    cd ../test
    go test -v main_test.go
