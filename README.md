# Consistent Hashing

A distributed key-value store built on consistent hashing with automatic Redis node health monitoring.

## How it works

- **Consistent hash ring**: Each Redis server is represented by 100 virtual nodes placed on a hash ring using Murmur3. Keys are assigned to the nearest virtual node clockwise.
- **Health monitoring**: A background goroutine pings all registered nodes every 3 seconds. Unhealthy nodes are skipped during reads/writes; the ring walks forward to the next healthy replica.
- **Failure recovery**: When a node comes back online, the coordinator re-establishes the Redis connection and marks it healthy, making it available for key operations again.

## Configuration

All configuration is via environment variables (defaults shown):

| Variable | Description | Default |
|---|---|---|
| `REDIS_URLS` | Comma-separated Redis node addresses | `localhost:6379, localhost:6380, localhost:6381, localhost:6382` |
| `VNODES_PER_NODE` | Virtual nodes per physical node on the hash ring | `100` |
| `PING_INTERVAL_SECONDS` | Seconds between health-check pings | `3` |
| `MAX_TRIES` | Max healthy nodes to try per operation (`0` = try all virtual nodes) | `0` |

## Quick start (Development)

```bash
docker compose up -d
```

Example `.env.local`:

```env
REDIS_URLS="localhost:6379, localhost:6380, localhost:6381, localhost:6382"
VNODES_PER_NODE=100
PING_INTERVAL_SECONDS=3
MAX_TRIES=0
```

```bash
go run ./cmd
```

## API

- `PUT /api/v1/keys/:key` — store a value (JSON body: `{"value": "...", "ttl": <seconds>}`) ttl = 0 for no ttl
- `GET /api/v1/keys/:key` — retrieve a value
- `GET /ping` — health check
- `GET /status` — node health state
