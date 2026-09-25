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

All endpoints return JSON. Success responses contain the data directly. Error responses use:

```json
{
  "error": "human-readable error message"
}
```

### Response Headers

| Header | Description |
|--------|-------------|
| `X-Cache` | `HIT` if served from cache, `MISS` if fetched from Bible Brain |
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
│   ├── api/                    # HTTP handlers (one file per endpoint)
│   │   ├── handler.go          # Shared Handler struct and JSON helpers
│   │   ├── languages.go        # GET /api/languages
│   │   ├── bibles.go           # GET /api/bibles
│   │   ├── books.go            # GET /api/books
│   │   ├── audio.go            # GET /api/audio
│   │   ├── copyright.go        # GET /api/copyright
│   │   └── health.go           # GET /health, GET /ready
│   ├── biblebrain/             # Bible Brain API client
│   │   ├── client.go           # HTTP client, URL builder
│   │   ├── bibles.go           # Fetch Bibles + copyright
│   │   ├── books.go            # Fetch book listings
│   │   ├── audio.go            # Fetch chapter audio
│   │   ├── gabon_config.go     # Gabon language/fileset allowlist
│   │   └── errors.go           # Typed error types
│   ├── cache/                  # Caching layer
│   │   ├── cache.go            # Cacher interface + factory
│   │   ├── memory.go           # In-memory implementation
│   │   └── redis.go            # Redis implementation
│   ├── config/                 # Configuration
│   │   └── config.go           # Environment variable loading
│   └── middleware/             # HTTP middleware
│       ├── cors.go             # CORS headers
│       ├── logging.go          # Structured request logging
│       └── security.go         # Security headers
├── Dockerfile                  # Multi-stage Docker build
├── docker-compose.yml          # Redis + server compose
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

### Cloud Run / Fly.io / Railway

The Docker image is self-contained. Set environment variables in your platform's dashboard:

1. Set `BIBLE_BRAIN_API_KEY` (required)
2. Set `REDIS_URL` if using managed Redis (recommended)
3. Set `CORS_ALLOWED_ORIGINS` to your Flutter app's domain
4. Set `ENV=production`
5. Deploy the Docker image

## Security

- The Bible Brain API key is **never exposed** to the Flutter client
- All incoming parameters are validated and sanitized
- The Gabon language/fileset allowlist is enforced **server-side**
- Security headers are set on all responses
- CORS is configurable per environment

## License

Private — All rights reserved.
