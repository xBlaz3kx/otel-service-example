# Service example with OpenTelemetry

This is an example HTTP service written in Go with OpenTelemetry support. Also included is configuration of the Grafana
stack (Grafana, Tempo, Loki, Prometheus, Promtail). It was made to learn how to make our services observable.

Disclaimer: The way some parts of the program are made might not be made properly or use best practices, since my
knowledge of observability and OpenTelemetry is limited.

## Running the example service

1. You can run the example in Docker by building the service:

```bash
docker network create grafana # create an external network
docker-compose build observability-example
docker-compose up -d
```

This will simplify the configuration.

2. Running on the host

```bash
go run cmd/app
```

## Running the Grafana stack

Simply run `docker-compose up -d`, and all the services should be up. You will need to manually add datasources (Loki,
Tempo and Prometheus) and enable linking between traces, metrics and logs. 