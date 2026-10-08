# Agnambie Backend

> Go backend proxy for the Agnambie audio Bible application — a Gabon-focused, offline-friendly Bible app.

## Architecture

```
Flutter App  →  Agnambie Go Backend  →  Bible Brain API
```

The Flutter application **never communicates directly** with the Bible Brain API. This backend:

- Securely holds the Bible Brain API key
- Proxies only the required Bible Brain endpoints
- Enforces the **Gabon language allowlist** server-side
- Normalises responses for the Flutter client
- Caches responses to respect Bible Brain rate limits

## Supported Languages

All Gabonese languages with audio Bible content available on Bible Brain:

| Code | Language | Filesets |
|------|----------|---------|
| FAN | Fang | NT (drama) |
| MYE | Myène | NT (audio + drama), OT |
| FRA | Français | NT + OT (multiple translations) |
| PUU | Punu | NT (audio + drama), OT |
| NZB | Nzebi / Njebi | NT (drama), OT |
| TSV | Tsogo | NT (audio + drama), OT |
| VIF | Vili | NT (audio + drama), portions |
| BKW | Bekwel | Portions |
| BBG | Barama | Portions |
| BUW | Bubi | Portions |
| DMA | Duma | Portions |
| KEB | Kélé | Portions |
| LUP | Lumbu | Portions |
| ZMN | Mbangwe | Portions |
| NMD | Ndumu | Portions |
| PIC | Pinji | Portions |
| SYX | Samay | Portions |
| SYI | Seki | Portions |
| BNG | Benga | Portions |

## API Endpoints

| Method | Path | Description | Cache TTL |
|--------|------|-------------|-----------|
| `GET` | `/health` | Liveness probe | — |
| `GET` | `/ready` | Readiness probe (verifies cache) | — |
| `GET` | `/api/languages` | List all Gabon languages and their filesets | — (static) |
| `GET` | `/api/bibles?language_code=FAN` | List Bibles for a language | 1 hour |
| `GET` | `/api/books?bible_id=FANBSG` | List books in a Bible | 1 hour |
| `GET` | `/api/audio?fileset_id=X&book=MAT&chapter=1` | Get chapter audio URL(s) | 30 min |
| `GET` | `/api/copyright?bible_id=FANBSG` | Get Bible copyright info | 24 hours |

### Query Parameters

#### `/api/bibles`
| Param | Required | Description |
|-------|----------|-------------|
| `language_code` | Yes | ISO 639-3 language code (e.g. `FAN`, `MYE`, `FRA`) |

#### `/api/books`
| Param | Required | Description |
|-------|----------|-------------|
| `bible_id` | Yes | Bible abbreviation (e.g. `FANBSG`) |

#### `/api/audio`
| Param | Required | Description |
|-------|----------|-------------|
| `fileset_id` | Yes | Audio fileset ID (e.g. `FANBSGN2DA`) |
| `book` | Yes | USFM book ID (e.g. `MAT`, `GEN`) |
| `chapter` | Yes | Chapter number (positive integer) |
| `quality` | No | Set to `data_saver` for opus16 low-bandwidth variant |

#### `/api/copyright`
| Param | Required | Description |
|-------|----------|-------------|
| `bible_id` | Yes | Bible abbreviation (e.g. `FANBSG`) |

### Response Format

All endpoints return JSON. Success responses are wrapped in a standard envelope:

```json
{
  "data": [ ... ],
  "meta": { "count": 1 }
}
```

Error responses use a structured format:

```json
{
  "error": {
    "code": 404,
    "message": "human-readable error message"
  }
}
```

### Advanced Low-Bandwidth Features (Gabon Optimizations)

To support users on constrained networks, this backend automatically implements:

- **Gzip Compression**: All JSON responses are compressed, reducing payload sizes by ~70-80%.
- **ETag / 304 Not Modified**: Endpoints generate weak ETags. If the Flutter client sends an `If-None-Match` header and the data hasn't changed, the backend returns a lightweight `304 Not Modified`.
- **Client Save-Data Header**: If the Android device sends the standard `Save-Data: on` HTTP header, the `/api/audio` endpoint automatically overrides the quality to the smallest `opus16` fileset when available.

### Response Headers

| Header | Description |
|--------|-------------|
| `X-Cache` | `HIT` if served from cache, `MISS` if fetched from Bible Brain |
| `ETag` | Hash of the response data (for conditional caching) |
| `Cache-Control` | Dictates how long the client can cache the response |
| `X-Content-Type-Options` | `nosniff` |
| `X-Frame-Options` | `DENY` |

## Getting Started

### Prerequisites

- **Go 1.27+**
- **Redis** (optional — an in-memory cache is used if `REDIS_URL` is not set)

### Setup

```bash
# Clone the repository
git clone https://github.com/adriel-meb/agnambie-backend.git
cd agnambie-backend

# Copy environment config
cp .env.example .env

# Edit .env and set your Bible Brain API key
# BIBLE_BRAIN_API_KEY=your-production-key-here

# Run locally (uses in-memory cache by default)
make run
```

### With Redis (recommended for production)

```bash
# Start Redis via Docker Compose
make docker-up

# Set REDIS_URL in .env
# REDIS_URL=redis://localhost:6379

# Run
make run
```

## Configuration

All configuration is via environment variables (12-factor app):

| Variable | Default | Description |
|----------|---------|-------------|
| `BIBLE_BRAIN_API_KEY` | — (required) | Bible Brain API key |
| `BIBLE_BRAIN_BASE_URL` | `https://4.dbt.io/api` | Bible Brain API base URL |
| `PORT` | `8080` | HTTP server port |
| `ENV` | `development` | `development` or `production` |
| `REDIS_URL` | — (optional) | Redis connection URL; if empty, in-memory cache is used |
| `CORS_ALLOWED_ORIGINS` | `*` | Comma-separated allowed origins |
| `LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |

## Logging

This project uses idiomatic, structured logging tailored for both development and production. 

- **Access Logs**: Every request generates exactly one `INFO` log detailing `method`, `path`, `status`, `duration`, `bytes` (response size), and `req_id`.
- **Application Logs**: Backend errors provide context-rich properties (e.g., `bible_id`, `fileset_id`) alongside the specific error.
- **Upstream Logs**: Communication with the Bible Brain API is tracked. Errors scrub any API keys from URLs to prevent leaks.
- **Request Tracing**: A unique `req_id` is generated for every request (or extracted from `X-Request-Id`) and is attached to all logs executed within that request lifecycle, allowing end-to-end tracing.

### Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Minimum level to log. Use `debug` to see upstream Bible Brain requests. |
| `LOG_FORMAT` | `text` | `text` for readable local development; `json` for machine-readable production logs. |

### Enabling Debug Logs

To troubleshoot upstream Bible Brain API issues, set `LOG_LEVEL=debug` locally:

```bash
LOG_LEVEL=debug LOG_FORMAT=text make run
```

### Examples

**Local Development (`text` format)**
```
time=2026-09-25T15:43:00.000Z level=INFO msg="server starting" port=8080 env=development
time=2026-09-25T15:43:02.105Z level=DEBUG msg="bible brain request completed" component=biblebrain upstream_path=/bibles/filesets/FANBSG/MAT/1 upstream_status=200 duration=150ms request_id=req-123
time=2026-09-25T15:43:02.106Z level=INFO msg="request completed" method=GET path=/api/audio query="fileset_id=FANBSG&book=MAT&chapter=1" status=200 bytes=421 duration=155ms remote=127.0.0.1 req_id=req-123
```

**Production (`json` format)**
```json
{"time":"2026-09-25T15:43:02.106Z","level":"ERROR","msg":"failed to fetch audio","fileset_id":"FANBSG","book":"MAT","chapter":1,"error":"upstream request failed: 404 Not Found","request_id":"req-123"}
{"time":"2026-09-25T15:43:02.106Z","level":"INFO","msg":"request completed","method":"GET","path":"/api/audio","query":"fileset_id=FANBSG&book=MAT&chapter=1","status":502,"bytes":78,"duration":"155ms","remote":"10.0.0.5","req_id":"req-123"}
```

## Development

```bash
# Run the server
make run

# Run all tests
make test

# Run tests with coverage
make cover

# Build binary
make build

# Lint (requires golangci-lint)
make lint

# Explore Bible Brain API (discovery tool)
make explore
```

## Project Structure

```
agnambie-backend/
├── cmd/
│   ├── server/main.go          # Application entrypoint
│   └── explore/main.go         # Bible Brain API explorer tool
├── internal/
│   ├── biblebrain/             # Bible Brain API client
│   │   ├── client.go           # HTTP client, URL builder
│   │   ├── bibles.go           # Fetch Bibles + copyright
│   │   ├── books.go            # Fetch book listings
│   │   ├── audio.go            # Fetch chapter audio
│   │   └── errors.go           # Typed error types
│   ├── cache/                  # Caching layer (Memory & Redis w/ singleflight)
│   ├── config/                 # Environment variable loading
│   ├── domain/                 # Core entities (Language, Bible, Audio) & Allowlist
│   ├── handler/                # HTTP request handlers & JSON envelopes
│   ├── middleware/             # HTTP middleware (CORS, Logging, Compress)
│   └── router/                 # Chi Router registration & middleware stack
├── Dockerfile                  # Multi-stage Docker build (runs as non-root)
├── docker-compose.yml          # Redis + server compose w/ healthchecks
├── Makefile                    # Build/run/test commands
├── .env.example                # Environment template
└── bible-brain-api-reference.md # Bible Brain API docs
```

## Deployment

### Docker

```bash
# Build and run everything
docker compose up --build

# Or build just the image
docker build -t agnambie-backend .
docker run -p 8080:8080 --env-file .env agnambie-backend
```

### Cloud Run / Fly.io / Railway / Render

The Docker image is self-contained. Set environment variables in your platform's dashboard:

1. Set `BIBLE_BRAIN_API_KEY` (required)
2. Set `REDIS_URL` if using managed Redis (recommended)
3. Set `CORS_ALLOWED_ORIGINS` to your Flutter app's domain
4. Set `ENV=production`
5. Deploy the Docker image

### Render (Free Tier) + Upstash Redis Deployment Guide

This is the recommended path for hosting the backend completely for free without requiring a credit card.

**1. Create the Database (Upstash)**
1. Go to [Upstash](https://upstash.com/) and create a free Redis database.
2. Select a US region (e.g., `us-east-1`) to minimize latency with Render.
3. In your database dashboard, under the **Connect** section, select **Node.js** or **ioredis** to reveal the native Redis connection string.
4. Copy the URL. It will look like this: `rediss://default:PASSWORD@your-endpoint.upstash.io:6379`. (Note the `rediss://` which enables required TLS).

**2. Deploy the App (Render)**
1. Sign up for [Render.com](https://render.com/) using GitHub.
2. Create a new **Web Service** and select **Build and deploy from a Git repository**.
3. Connect the `agnambie-backend` repository.
4. Set the **Root Directory** to `agnambie-backend`.
5. Render will automatically detect the `Dockerfile` and select the **Docker** environment.
6. Ensure the **Free** instance type is selected.
7. Add the following Environment Variables:
   - `REDIS_URL`: Paste the `rediss://` URL from Upstash.
   - `BIBLE_BRAIN_API_KEY`: Your real API key.
8. Click **Create Web Service**. 

**3. Mitigate "Freezing" (The Keep-Alive Ping)**
Render puts Free Web Services to sleep after 15 minutes of inactivity. When it wakes up, the first request takes ~10-30 seconds. To keep the API instant 24/7 without exceeding the 750 free hours/month limit:
1. Create a free account on [UptimeRobot.com](https://uptimerobot.com/) or [cron-job.org](https://cron-job.org/).
2. Create an HTTP monitor pointing to your Render app's health endpoint: `https://your-app.onrender.com/health`.
3. Set the ping interval to **14 minutes**.
4. The Go backend processes this ping with near-zero resources, effectively tricking Render into keeping the container awake 24/7.

## Security

- The Bible Brain API key is **never exposed** to the Flutter client
- All incoming parameters are validated and sanitized
- The Gabon language/fileset allowlist is enforced **server-side**
- Security headers are set on all responses
- CORS is configurable per environment

## License

Private — All rights reserved.
