# async-api

Async API example application.

> AsyncAPI is an open source initiative that seeks to improve the current state of Event-Driven Architectures (EDA).

Based on <a href="https://github.com/morzhanov/go-otel">go-otel</a> repo.

<img src="https://i.ibb.co/gm81Csc/Untitled-2021-10-18-1452.png" alt="arch"/>

## App Description

Simple app: Orders and Payments services + API GW.

- API GW handles REST requests and proxies them to services
- Orders - create and process orders
- Payments - process payment and return payment info
- Orders handles REST requests from API GW
- Payments handles gRPC requests from API GW
- Payments handles Kafka events from API GW
- Orders saves data to MongoDB
- Payments saves data to PostgreSQL

## AsyncAPI schema

Application uses AsyncAPI generated schema which could be used for UI schema visualization (like OpenAPI for REST services).

In order to use AsyncAPI schema:
- run each service
- `yml` schema files should be generated in the `/api/*` directories
- run AsyncAPI UI client and use `yml` files for schame visualisation

## Structure

- `/api` - contains Kafka, REST and gRPC service and message definitions
- `/cmd` - application setup
- `/config` - .env file with environment variables
- `/deploy`
    - `docker-compose.yml` - docker-compose file with Jaeger, Prometheus, MongoDB and PostgreSQL setup
- `/internal`
    - `/apigw` - API GW service internals
    - `/config` - config files setup with viper
    - `/event` - events base controller
    - `/grpc` - grpc base controller
    - `/logger` - application logger, creates file transport (for filebeat) and console transport
    - `/mongodb` - mongodb database setup
    - `/order` - order service internals
    - `/payment` - payment service internals
    - `/psql` - postgres database setup
    - `/rest` - application REST base controller
    - `/telemetry` - otel setup files

## Local Running

You should deploy dependencies with docker-compose:

```
cd deploy
docker-compose up -d
```

Then you could run all services separately:

```
go run ./cmd/<service-name>/main.go
```

