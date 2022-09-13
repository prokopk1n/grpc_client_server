# Клиент-серверное приложение с использованием фреймворка gRPC.

# Запуск приложения:
    go run server/main.go
    cd client/src
    go run main.go
# Тестирование приложения:
    cd docker
    docker-compose -p newtest -f docker-compose-test.yml -d up
    cd ../test
    go test -v main_test.go

## Table of contents
* [General info](#general-info)
* [Technologies](#technologies)
* [Setup](#setup)
