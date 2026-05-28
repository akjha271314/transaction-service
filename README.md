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

## Swagger UI

Once the server is running, open:

```
http://localhost:8080/swagger/index.html
```

The Swagger UI is available without authentication. Use the **Authorize** button to set your `X-API-Key` before trying out the protected endpoints.

To regenerate the docs after changing handler annotations:

```bash
swag init --generalInfo main.go --output docs
```

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
{ "document_number": "12345678900", "balance": 500.0 }
```

`balance` is the initial available funds and defaults to `0`. Must be `>= 0`. Balance changes with every transaction — it is not a static limit.

**Response** `201 Created`
```json
{ "account_id": 1, "document_number": "12345678900", "balance": 500.0 }
```

---

### Get an account

```
GET /accounts/:accountId
```

**Response** `200 OK`
```json
{ "account_id": 1, "document_number": "12345678900", "balance": 460.0 }
```

`balance` reflects the current available funds at the time of the request.

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

Returns `422` if the account or operation type does not exist, or if the transaction would exceed the account's credit limit.

---

## Operation types

| ID | Description                | `is_credit` |
|----|----------------------------|-------------|
| 1  | Normal Purchase            | false       |
| 2  | Purchase with installments | false       |
| 3  | Withdrawal                 | false       |
| 4  | Credit Voucher             | true        |

The `is_credit` flag on each operation type determines the sign of the stored amount — no hardcoded IDs in business logic. Adding a new operation type is a data change only: insert a row into `operation_types` with the correct `is_credit` value.

Each transaction atomically updates the account `balance`. A transaction is rejected with `422` if it would bring the balance below `0`. To add funds, create a credit voucher transaction.

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
