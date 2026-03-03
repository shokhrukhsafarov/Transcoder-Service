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

## How Transcoding Works

Transcoding is orchestrated entirely through RabbitMQ. This service acts as the **control plane** — it does not run FFmpeg itself; an external transcoder worker consumes the job and does the actual encoding.

### Step-by-step flow

```
Client → POST /pipeline → REST API → PostgreSQL → RabbitMQ (write queue)
                                                         ↓
                                               External transcoder worker
                                                         ↓
                                          RabbitMQ (listen queue) → REST API
                                                         ↓
                                             PostgreSQL updated + webhook fired
```

**1. Job submission (`POST /pipeline`)**

The client sends a pipeline creation request with:
- `input_url` — publicly reachable URL of the source video
- `output_key` — destination filename / HLS manifest key in the CDN bucket
- `output_path` — CDN bucket name
- `project_id` — links the job to a project and its associated CDN storage credentials
- Optional: `resolutions`, `audio_tracks`, `subtitle`, `drm`, `key_id`

The handler fetches the file size from `input_url`, persists the pipeline record in PostgreSQL (status: pending), then publishes a `PipelineRabbitMq` message to the **write queue** (`WRITE_QUEUE`). The message includes all CDN credentials (access key, secret, region, bucket, domain) and the DRM key ID so the worker can operate autonomously.

**2. External transcoder worker**

The external worker (not part of this repo) consumes messages from the write queue and performs three sequential stages:
- **preparation** — download/probe source file, extract metadata
- **transcode** — encode to the requested resolutions using FFmpeg (or similar); supports multi-audio-track muxing, subtitle embedding, and DRM encryption
- **upload** — push output segments/manifests to the configured CDN bucket (S3-compatible)

Each stage is timed independently (preparation / transcode / upload duration in ms).

**3. Status updates (listen queue)**

When a stage completes or fails the worker publishes an `UpdatePipelineStatus` message back to the **listen queue** (`LISTEN_QUEUE`). The service consumes these messages in `StartListening()` and updates the pipeline row with:
- `stage` — which stage just finished (`preparation`, `transcode`, `upload`)
- `stage_status` — `success` or `fail`
- Per-stage duration in seconds
- Produced resolutions (resolution, measure, bitrate)
- `fail_description` on errors

**4. Webhook notification**

Once the job either fails at any stage or reaches `upload → success`, the service looks up the active webhook registered for the project and fires a `POST` request to the webhook URL, signed with the project's access/secret key pair. The pipeline row records `webhook_status` (boolean) and `webhook_retry_count`.

**5. Webhook retry**

`Retry()` runs on a schedule and retries any pipeline whose webhook call was unsuccessful (up to the configured maximum), querying for pipelines where `webhook_status = false`.

### Pipeline stages

| Stage | Description |
|---|---|
| `preparation` | Download and probe the source file |
| `transcode` | Encode to target resolutions, mux audio tracks, embed subtitles, apply DRM |
| `upload` | Push HLS/DASH segments and manifest to CDN |

### DRM support

Setting `drm: true` and supplying a `key_id` in the pipeline request instructs the worker to apply encryption during transcoding. The key ID is forwarded as-is in the RabbitMQ message.

---

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
