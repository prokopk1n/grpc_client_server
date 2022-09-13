# Клиент-серверное приложение с использованием фреймворка gRPC.
Приложение предусмотрено для хранения авиабилетов или поиска информации по ним. 
![Image alt](https://github.com/prokopk1n/resources/raw/grpc_client_server/main_page.png)

## Table of Contents
* [Setup](#setup)
* [General info](#general-info)
* [Technologies](#technologies)
* [More info](#more-info)

## Setup
### Install and launch:
    git clone https://github.com/prokopk1n/grpc_client_server
    docker-compose -p db -f docker-compose-db.yml up
    docker-compose -p app -f docker-compose.yml up
### Launch tests:
    cd docker
    docker-compose -p newtest -f docker-compose-test.yml -d up
    cd ../test
    go test -v main_test.go

## General info
За основу взята тестовая база данных https://postgrespro.ru/education/demodb <br /> 
Доступ к сайту после запуска осуществляется по адресу https://localhost:8080/ <br /><br />
Для работы с базой данных билетов написан отдельный сервис, который можно найти в директории server. <br /><br />
Взаимодействие с веб-клиентом осуществляет другой сервис, код которого можно найти в директории client. Доступ к базе данных билетов происходит с помощью фреймфорка gRPC. </br>
Для аутентификации пользователей используется JSON Web Token. Для хранения refresh token, а также логинов и хэш-сумм паролей пользователей создана отдельная база данных clientdb.

## Technologies
В проекте используются:
* golang version: 1.19
* postgresql version: 10.5
* grpc-go release: 1.49.0
* docker release: 20.10.16

## More info
### Схема базы данных clientdb:
![Image alt](https://github.com/prokopk1n/resources/raw/grpc_client_server/clientdb.png)
