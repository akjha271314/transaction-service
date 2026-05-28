### Transaction Service

A REST API service for managing cardholder accounts and their transactions.

---

## Requirements

- Go 1.22+
- (Optional) Docker & Docker Compose

---

## Configuration

All configuration is read from environment variables. Copy `.env.example` to `.env` and set the values before running locally.

```bash
cp .env.example .env
```

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `API_KEY` | Yes | — | Secret key sent by clients in the `X-API-Key` header |
| `PORT` | No | `8080` | Port the server listens on |

In production, set these variables directly in your environment — no `.env` file is needed.

---

## Running the service

### Option 1 — Go directly

```bash
./run
```

Or equivalently:

```bash
go run .
```

### Option 2 — Docker Compose

```bash
docker compose up --build
```

The server starts on the port defined by `PORT` (default `8080`).

---

## Running tests

```bash
go test ./...
```

---

## API

All endpoints except `/health` require an `X-API-Key` header:

```
X-API-Key: your-secret-api-key
```

Returns `401 Unauthorized` if the header is missing or incorrect.

---

### Health check

```
GET /health
```

**Response** `200 OK` — no authentication required.

---

### Create an account

```
POST /accounts
```

**Request**
```json
{ "document_number": "12345678900" }
```

**Response** `201 Created`
```json
{ "account_id": 1, "document_number": "12345678900" }
```

---

### Get an account

```
GET /accounts/:accountId
```

**Response** `200 OK`
```json
{ "account_id": 1, "document_number": "12345678900" }
```

Returns `404` if the account does not exist.

---

### Create a transaction

```
POST /transactions
```

**Request**
```json
{ "account_id": 1, "operation_type_id": 4, "amount": 123.45 }
```

**Response** `201 Created`
```json
{
  "transaction_id": 1,
  "account_id": 1,
  "operation_type_id": 4,
  "amount": 123.45,
  "event_date": "2024-01-05T09:34:18Z"
}
```

Returns `422` if the account or operation type does not exist.

---

## Operation types

| ID | Description                           |
|----|---------------------------------------|
| 1  | Normal Purchase (negative amount)     |
| 2  | Purchase with installments (negative) |
| 3  | Withdrawal (negative)                 |
| 4  | Credit Voucher (positive)             |

Purchases and withdrawals are always stored with **negative** amounts; credit vouchers with **positive** amounts — regardless of the sign sent in the request.

---

## Architecture

```
main.go
internal/
  config/      — Environment variable loading
  db/          — SQLite in-memory setup & migrations
  models/      — Account, Transaction structs
  repository/  — Database access layer
  service/     — Business logic (amount sign enforcement, validation)
  handler/     — HTTP request/response layer
  router/      — Route definitions and API key middleware
  testutil/    — Shared test helpers
```

Data is stored in-memory (SQLite `:memory:`) and resets on each restart.
