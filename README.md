# Shipping

This is an example of real world production ready application consisting of multiple microservices written in different programming languages.

> The implementation is based on the container shipping domain from the [Domain Driven Design](http://www.amazon.com/Domain-Driven-Design-Tackling-Complexity-Software/dp/0321125215) book by Eric Evans, which was [originally](http://dddsample.sourceforge.net/) implemented in Java but has since been ported to Go. The [original Go application](https://github.com/marcusolsson/goddd) is maintained separately and comes with a mock [routing service](https://github.com/marcusolsson/pathfinder).

## Motivation

This project aims to demonstrate a possible way to split a monolithic application into microservices and organize synchronous and asynchronous communication between them. This project is **not** a DDD development best practice or idiomatic way of developing in Go and Rust.

## Architecture Overview

The application consists of four application services and several other services for communication and maintenance. [Protobuf](https://developers.google.com/protocol-buffers/) is chosen as the IDL for synchronous communication between services via [gRPC](https://grpc.io/) and as the data format for asynchronous communication via [RabbitMQ](https://www.rabbitmq.com/).

### App features

* **booking** - used by the shipping company to book and route cargos.
* **handling** - used by our staff around the world to register whenever the cargo has been received, loaded etc.
* **tracking** - used by the customer to track the cargo along the route.

### App implementation

* **proto** is a folder contains API definitions of all services and describes the data format for RabbitMQ messages.
* **booking** is a booking service written in Go.
* **pathfinder** is a mock routing service wrapped in gRPC and used by booking service.
* **handling** is a handling service written in Rust.
* **tracking** is a tracking service written in Go.
* **apigateway** is an [Envoy](https://www.envoyproxy.io/docs/envoy/latest/start/install) that is used to transcode gRPC-JSON to provide API for browser clients.
* **frontend** is a use case written in Rust with [Seed](https://seed-rs.org/).
* **elk** is an [ELK stack](https://www.elastic.co/downloads/) for centralized logging and maintenance.

## Running the application

Make sure you have installed [docker](https://docs.docker.com/engine/install/), [docker-compose](https://docs.docker.com/compose/), [make](https://www.gnu.org/software/make/), [protoc](https://grpc.io/docs/protoc-installation/) and [Go plugins for protoc](https://grpc.io/docs/languages/go/quickstart/). After that, you can run the below command from the project directory.

```bash
$ make run
```

After everything has bootstrapped, navigate to http://localhost/. You can watch the original [video](https://youtu.be/eA8xgdtqqs8) with use cases.
