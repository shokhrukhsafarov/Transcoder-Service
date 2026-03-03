# Transcoder Service

A microservice for managing video transcoding pipelines. It provides a REST API to create and track transcoding jobs organized around companies, projects, and pipelines. Pipeline status updates are consumed from a message queue, and webhook callbacks notify external systems on job completion.

## Stack

| Layer | Technology |
|---|---|
| Language | Go 1.23.2 |
| HTTP framework | Gin |
| Database | PostgreSQL + sqlx + Squirrel |
| Cache | Redis |
| Message broker | RabbitMQ |
| Auth | JWT + Casbin (RBAC) |
| RPC | gRPC + Protobuf |
| API docs | Swagger (swag) |
| Logging | zerolog |
| DI | uber/fx |
| Containerization | Docker / Docker Compose |
| Orchestration | Kubernetes (Helm) |

## Getting Started

### Prerequisites

- Go 1.23.2+
- Docker & Docker Compose
- Make

### Setup

```bash
cp .env.dist .env
# Fill in your values in .env
```

### Run

```bash
make compose_up   # start all services
make compose_down # stop all services
```

Once running:

- REST API + Swagger: `http://localhost:8000/v1/swagger/index.html`
- gRPC: `localhost:9110`
- RabbitMQ management: `http://localhost:15672`

### Useful commands

```bash
make run            # run without Docker
make build          # build binary
make swag_init      # regenerate Swagger docs
make migrate_up     # run DB migrations
make gen-proto-module # regenerate gRPC code from proto files
make linter         # run linter
```

## Environment Variables

See `.env.dist` for the full list. Key variables:

| Variable | Description |
|---|---|
| `POSTGRES_HOST` | PostgreSQL host |
| `POSTGRES_PASSWORD` | PostgreSQL password |
| `RABBITMQ_HOST` | RabbitMQ host |
| `RABBITMQ_USER` / `RABBITMQ_PASSWORD` | RabbitMQ credentials |
| `SIGN_IN_KEY` | JWT signing key |
| `LISTEN_QUEUE` | Queue to consume pipeline status updates from |
| `WRITE_QUEUE` | Queue to publish pipeline jobs to |
