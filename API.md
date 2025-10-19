# API Documentation

**Base URL**: `http://localhost:8080`

## Endpoints

### Health Check
`GET /healthz`

```bash
curl http://localhost:8080/healthz
```

### Version
`GET /version`

```bash
curl http://localhost:8080/version
```

### Solve Packs
`POST /packs/solve`

Find optimal pack combination for given amount.

**Request:**
```bash
curl -X POST http://localhost:8080/packs/solve \
  -H "Content-Type: application/json" \
  -d '{
    "sizes": [250, 500, 1000, 2000, 5000],
    "amount": 12001
  }'
```

**Response:**
```json
{
  "solution": {
    "250": 1,
    "5000": 2
  },
  "overage": 249,
  "packs": 3
}
```

**Validation:**
- `sizes`: array > 0, values ≤ 1,000,000
- `amount`: > 0 and ≤ 1,000,000,000

**Status Codes:**
- `200` - success
- `400` - invalid JSON
- `422` - validation error
- `500` - internal error

## Features

- **Correlation ID**: Optional `X-Correlation-ID` header for request tracing
- **Idempotency**: Identical requests return identical results
- **Structured Logging**: JSON logs with correlation ID
- **Graceful Shutdown**: Clean shutdown on SIGINT/SIGTERM
