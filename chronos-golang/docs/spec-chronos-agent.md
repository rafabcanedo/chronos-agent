# Personal Finance App — Backend & Supabase Technical Specification

## Overview
This document specifies the backend system and database architecture for a lightweight personal finance tracking application. The backend is written in **Go** and exposes a RESTful JSON API. Data persistence is handled by **PostgreSQL hosted on Supabase**. Natural language interactions are parsed using **Abacus RouteLLM** via an OpenAI-compatible HTTP interface.

## 1. Technology Stack & Infrastructure

- **Language:** Go (modern standard toolchain)
- **Database:** PostgreSQL (Supabase Free Tier)
- **Database Driver:** `pgx` (recommended) or standard `database/sql` + `lib/pq`
- **Authentication:** Stateless JWT HS256 (Bearer Token) via `obsidian-golang-toolkit/api/auth`
- **Password Hashing:** `bcrypt` (external dependency — not covered by the toolkit)
- **AI Integration:** Abacus RouteLLM API (`https://routellm.abacus.ai/v1/chat/completions`)
- **Hosting Target:** Containerized service on free-tier cloud host (e.g., Render, Fly.io, Koyeb)
- **Internal Toolkit:** `github.com/rafabcanedo/obsidian-golang-toolkit` (see Section 1.1)

## 1.1 Obsidian Golang Toolkit

The Obsidian Golang Toolkit (`github.com/rafabcanedo/obsidian-golang-toolkit`) is a personal collection of framework-agnostic Go utility packages. It has no external dependencies beyond the standard library (except `gopkg.in/yaml.v3` for YAML config loading). It is imported selectively — only the packages needed by this project.

### Packages used in this project

| Toolkit package | Replaces / covers | Used in |
|---|---|---|
| `api/auth` | `internal/auth` entirely | JWT generation (`Generate`), parsing (`Parse`), and HTTP middleware (`Middleware`, `GetClaims`) |
| `api/httpclient` | Raw `net/http` client in `internal/llm` | Calling Abacus RouteLLM API with timeout, Bearer header, and retry |
| `api/middleware` | `internal/httpx` middleware | `Logger`, `Recovery`, `RequestID`, `CORS` applied in `main.go` |
| `api/respond` | `internal/httpx` response helpers | `OK`, `Created`, `NoContent`, `Error`, `Paginated` used in all handlers |
| `api/binding` | Manual `json.Decode` in handlers | Decoding JSON request bodies and query params into structs |
| `config` | `internal/config` env loading | Loading `.env` file on startup |
| `env` | `os.Getenv` calls | Reading individual env vars with defaults (`DATABASE_URL`, `JWT_SECRET`, `ABACUS_API_KEY`) |
| `util/errutil` | Ad-hoc error strings | Typed HTTP errors (`NotFound`, `BadRequest`, `Unauthorized`) propagated from services to handlers |
| `util/logutil` | `log.Printf` calls | Structured logging via `slog` throughout the application |
| `util/contextutil` | Manual context key types | Used internally by `api/auth` and `api/middleware`; also available for passing `user_id` in context |
| `validator` | Manual field checks in handlers | Fluent input validation (`Required`, `MinLen`, `MaxLen`) on request bodies |
| `health` | No equivalent in spec | Adds a `GET /health` endpoint for uptime monitoring |

### Note on `internal/auth` and `internal/httpx`

Because `api/auth` and `api/middleware` + `api/respond` from the toolkit cover the full responsibility of these two internal packages, `internal/auth` and `internal/httpx` will be thin wrappers or may be removed entirely. The toolkit packages are used directly in handlers and `main.go`.

## 2. Database Schema (PostgreSQL on Supabase)

### 2.1 Table: `users`
Stores credentials for application access.

id
UUID
PRIMARY KEY DEFAULT gen_random_uuid()
Unique user identifier

username
VARCHAR(50)
NOT NULL UNIQUE
Login handle

email
VARCHAR(255)
NOT NULL UNIQUE
User email address

password
TEXT
NOT NULL
bcrypt hash of the user password

created_at
TIMESTAMPTZ
NOT NULL DEFAULT NOW()
Record creation timestamp

Note: User registration is handled via the `POST /api/auth/user` endpoint.

### 2.2 Table: `recurring_bills`
Defines recurring monthly bill templates ("blueprints").

id
UUID
PRIMARY KEY DEFAULT gen_random_uuid()
Unique template identifier

user_id
UUID
NOT NULL REFERENCES users(id) ON DELETE CASCADE
Owner identifier

name
VARCHAR(100)
NOT NULL
Description (e.g., "Internet", "Rent")

default_amount
NUMERIC(12,2)
NOT NULL CHECK (default_amount >= 0)
Expected default bill amount

billing_day
SMALLINT
NOT NULL CHECK (billing_day BETWEEN 1 AND 31)
Day of the month the bill is due

category
VARCHAR(50)
NOT NULL DEFAULT 'outros'
Expense category — must match a value from the shared category enum defined in the backend and frontend: 'mercado', 'restaurantes', 'moradia', 'transporte', 'lazer', 'outros'

active
BOOLEAN
NOT NULL DEFAULT TRUE
Whether to generate instances each month

created_at
TIMESTAMPTZ
NOT NULL DEFAULT NOW()
Record creation timestamp

updated_at
TIMESTAMPTZ
NOT NULL DEFAULT NOW()
Last update timestamp — should be refreshed on every UPDATE via a trigger or explicitly in the repository layer

### 2.3 Table: `transactions`
Records actual ledger entries (both generated recurring bills and one-off variable expenses).

id
UUID
PRIMARY KEY DEFAULT gen_random_uuid()
Unique transaction identifier

user_id
UUID
NOT NULL REFERENCES users(id) ON DELETE CASCADE
Owner identifier

recurring_bill_id
UUID
NULL REFERENCES recurring_bills(id) ON DELETE SET NULL
Reference to recurring template; NULL for ad-hoc expenses

description
VARCHAR(150)
NOT NULL
Transaction description

amount
NUMERIC(12,2)
NOT NULL CHECK (amount >= 0)
Final transaction value

category
VARCHAR(50)
NOT NULL DEFAULT 'outros'
Category — same enum as recurring_bills: 'mercado', 'restaurantes', 'moradia', 'transporte', 'lazer', 'outros'

due_date
DATE
NULL
Due date; can be NULL for immediate one-off expenses

status
VARCHAR(20)
NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid'))
Current lifecycle state

paid_at
DATE
NULL
Date payment was executed; populated only when status = 'paid'

created_at
TIMESTAMPTZ
NOT NULL DEFAULT NOW()
Record creation timestamp

updated_at
TIMESTAMPTZ
NOT NULL DEFAULT NOW()
Last update timestamp — refreshed on every UPDATE

### 2.4 Integrity & Deduplication Rules
Monthly recurring deduplication: A unique compound index prevents generating duplicate transactions for the same recurring bill within the same calendar month:

```sql
CREATE UNIQUE INDEX idx_unique_recurring_per_month
ON transactions (user_id, recurring_bill_id, (DATE_TRUNC('month', due_date)))
WHERE recurring_bill_id IS NOT NULL;
```

Monetary values: All currencies strictly use NUMERIC(12,2) to prevent floating-point rounding inaccuracies.

Overdue state: There is no dedicated overdue database column. A record is computed as overdue by the application when `status = 'pending' AND due_date < CURRENT_DATE`.

## 3. Go Project Structure & Architectural Layers
The service follows a domain-driven, layered architecture without unnecessary enterprise abstractions.

```
backend/
├── cmd/
│   └── api/
│       └── main.go          # Dependency wiring, DB connection, HTTP server start
├── internal/
│   ├── config/              # Environment variable loading (DATABASE_URL, JWT_SECRET, ABACUS_API_KEY)
│   ├── database/            # Supabase Postgres connection pool initialization
│   ├── auth/                # JWT generation, token parsing, and auth middleware
│   ├── transactions/        # Handler, Service, Repository for financial entries
│   ├── recurring/           # Handler, Service, Repository for bill templates & monthly generation
│   ├── chat/                # Chat handler and business orchestration (LLM -> Service action)
│   ├── llm/                 # HTTP client calling Abacus RouteLLM API
│   └── httpx/               # Common HTTP responses, error formatters, CORS middleware
├── migrations/              # Plain .sql migration scripts
├── .env.example
├── Dockerfile
└── go.mod
```

**Layer Responsibilities**

- **Handlers (HTTP):** Decode JSON payloads, extract user context (`user_id` from JWT), enforce basic validation, return JSON responses.
- **Services (Business Logic):** Contain domain operations (e.g., fuzzy-matching a bill description, synchronizing `status = 'paid'` with `paid_at`, generating monthly bills).
- **Repositories (Data Access):** Execute raw or prepared SQL queries against Supabase. Enforce `WHERE user_id = $1` on every query.
- **Chat Orchestrator (`chat` package):** Calls the `llm` package to get a structured JSON action, performs validation, and delegates execution to the appropriate service.

## 4. API Endpoints Specification

All endpoints except `POST /api/auth/login` and `POST /api/auth/user` require the `Authorization: Bearer <JWT_TOKEN>` header.

All successful responses use a `{ "data": ... }` envelope. Errors return `{ "error": "message" }`.

---

### 4.1 Authentication

**`POST /api/auth/user`** — public

Creates a new user account.

```json
// request
{
  "username": "string",
  "email": "string",
  "password": "string"
}

// response 201
{
  "data": {
    "id": "uuid",
    "username": "string",
    "email": "string",
    "created_at": "timestamp"
  }
}

// response 400
{ "error": "username already taken" }
```

---

**`POST /api/auth/login`** — public

Authenticates a user and returns a JWT.

```json
// request
{
  "username": "string",
  "password": "string"
}

// response 200
{ "data": { "token": "jwt_string" } }

// response 401
{ "error": "invalid credentials" }
```

Token expiry: 30 days. Contains `user_id` in claims.

---

### 4.2 Transactions

**`GET /api/transactions?month=YYYY-MM&status=pending|paid`**

Fetches user transactions filtered by calendar month and/or payment state. `month` defaults to current month if omitted.

**Month filter logic:** Recurring transactions have a `due_date` and ad-hoc expenses have `due_date = NULL` but always have `paid_at` set. To include both types in a monthly query, the repository must filter using `COALESCE(due_date, paid_at)`:

```sql
WHERE user_id = $1
  AND DATE_TRUNC('month', COALESCE(due_date, paid_at)) = DATE_TRUNC('month', $2::date)
```

```json
// response 200
{
  "data": [
    {
      "id": "uuid",
      "recurring_bill_id": "uuid | null",
      "description": "string",
      "amount": 0.00,
      "category": "string",
      "due_date": "YYYY-MM-DD | null",
      "status": "pending | paid",
      "paid_at": "YYYY-MM-DD | null",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ]
}
```

---

**`PATCH /api/transactions/:id`**

Updates one or more mutable fields of a transaction. All fields are optional — only provided fields are updated.

Mutable fields: `description`, `amount`, `category`, `due_date`, `status`, `paid_at`.

```json
// request (all fields optional)
{
  "description": "string",
  "amount": 0.00,
  "category": "string",
  "due_date": "YYYY-MM-DD",
  "status": "pending | paid",
  "paid_at": "YYYY-MM-DD"
}

// response 200
{
  "data": {
    "id": "uuid",
    "recurring_bill_id": "uuid | null",
    "description": "string",
    "amount": 0.00,
    "category": "string",
    "due_date": "YYYY-MM-DD | null",
    "status": "pending | paid",
    "paid_at": "YYYY-MM-DD | null",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}

// response 404
{ "error": "transaction not found" }
```

---

**`DELETE /api/transactions/:id`**

Permanently removes a transaction record.

```json
// response 200
{ "data": { "message": "Transação removida com sucesso." } }

// response 404
{ "error": "transaction not found" }
```

---

### 4.3 Recurring Bills (Templates)

**`GET /api/recurring?active=true|false`**

Lists recurring bill templates. Defaults to `active=true`. Pass `active=false` to list inactive bills.

```json
// response 200
{
  "data": [
    {
      "id": "uuid",
      "name": "string",
      "default_amount": 0.00,
      "billing_day": 1,
      "category": "string",
      "active": true,
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ]
}
```

---

**`POST /api/recurring`**

Creates a new recurring bill template.

```json
// request
{
  "name": "string",
  "default_amount": 0.00,
  "billing_day": 1,
  "category": "string"
}

// response 201
{
  "data": {
    "id": "uuid",
    "name": "string",
    "default_amount": 0.00,
    "billing_day": 1,
    "category": "string",
    "active": true,
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}

// response 400
{ "error": "validation error message" }
```

---

**`PUT /api/recurring/:id`**

Updates an existing recurring bill template. All fields are optional. Changes apply only to future generated transactions — existing transactions are not retroactively modified. Set `active: false` to deactivate a bill (stops generating new monthly instances).

```json
// request (all fields optional)
{
  "name": "string",
  "default_amount": 0.00,
  "billing_day": 1,
  "category": "string",
  "active": false
}

// response 200
{
  "data": {
    "id": "uuid",
    "name": "string",
    "default_amount": 0.00,
    "billing_day": 1,
    "category": "string",
    "active": false,
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}

// response 404
{ "error": "recurring bill not found" }
```

---

### 4.4 Chat & NLP Processing

**`POST /api/chat`**

```json
// request
{ "message": "string" }
```

Flow:
1. Extracts `user_id` from JWT context.
2. Sends prompt and current server date (`YYYY-MM-DD`) to Abacus RouteLLM.
3. Parses structured JSON response from the LLM.
4. Executes the corresponding database operation.
5. Returns a PT-BR confirmation message or a descriptive error.

```json
// response 200
{ "data": { "message": "Conta de internet marcada como paga." } }

// response 404 (pay action — bill not found)
{ "error": "Conta não encontrada para este mês." }

// response 400 (unrecognized intent or missing fields)
{ "error": "Não entendi o que você quis dizer." }
```

---

## 5. Abacus RouteLLM Integration & Prompt Contract

### 5.0 Language Policy

**MVP (current version):** The agent communicates exclusively in **Brazilian Portuguese (PT-BR)**. All user messages sent to the chat endpoint are expected to be in PT-BR, and the system prompt examples reflect this. The LLM is expected to handle PT-BR input correctly since Abacus RouteLLM routes to capable multilingual models.

**Future versions:** English support is planned. When added, the system prompt will be updated to handle both languages, and the examples in the prompt contract below will be extended accordingly.

The backend uses standard HTTP POST to communicate with Abacus RouteLLM.

- Endpoint: `https://routellm.abacus.ai/v1/chat/completions`
- Model: `route-llm`
- Headers: `Authorization: Bearer <ABACUS_API_KEY>`, `Content-Type: application/json`

**Prompting Strategy**

The backend sets a strict system prompt instructing the model to output only raw JSON with no markdown fences or prose commentary:

```
You are an intent parser for a personal finance system.
Today's date is: {{CURRENT_DATE}}.

Given the user's message, classify the intent into one of the following actions:
- "pay": User wants to mark an existing bill as paid (e.g., "paguei a internet", "conta de luz paga").
- "create_expense": User incurred a one-off/ad-hoc expense (e.g., "gastei 20 no sorvete", "paguei 50 no almoço").
- "create_recurring": User wants to define a recurring bill (e.g., "cadastrar aluguel de 1200 todo dia 5").
- "query": User is asking a financial question (e.g., "quanto gastei esse mês?").

Respond ONLY with a valid JSON object matching this schema:
{
  "action": "pay" | "create_expense" | "create_recurring" | "query",
  "description": string,
  "amount": number or null,
  "date": "YYYY-MM-DD" or null,
  "billing_day": integer (1-31) or null,
  "category": string or null
}
```

**Action Handling Logic in Backend**

`action == "pay"`:
Go runs an ILIKE / pattern match for pending transactions in the current month matching description:

```sql
SELECT id FROM transactions
WHERE user_id = $1
  AND status = 'pending'
  AND description ILIKE '%' || $2 || '%'
  AND DATE_TRUNC('month', due_date) = DATE_TRUNC('month', $3::date)
LIMIT 1;
```

If found: Updates `status = 'paid'`, `paid_at = date`, returns success.
If not found: Does NOT auto-create a bill. Returns HTTP 404 with `"Conta não encontrada para este mês."`.

`action == "create_expense"`:
Inserts immediately into `transactions` with `status = 'paid'`, `paid_at = date`, `recurring_bill_id = NULL`. Returns success.

`action == "create_recurring"`:
Inserts into `recurring_bills` (`name`, `default_amount`, `billing_day`, `category`). Returns success.

`action == "query"`:
Backend performs aggregate query (`SUM(amount)` grouped by category or status) for the requested month, passes numbers back to the LLM or formats directly in Go, and returns the answer.

## 6. Monthly Recurring Bills Generation (Backend Logic)
To ensure bills for the new month exist without needing external cron jobs, the backend uses lazy evaluation:

When a user requests `GET /api/transactions` or `POST /api/chat`, the backend checks if bills for the current month (`DATE_TRUNC('month', NOW())`) have already been populated for this `user_id`.

If unpopulated, it executes an idempotent batch INSERT:

```sql
INSERT INTO transactions (user_id, recurring_bill_id, description, amount, category, due_date, status)
SELECT
    rb.user_id,
    rb.id,
    rb.name,
    rb.default_amount,
    rb.category,
    MAKE_DATE(
        EXTRACT(YEAR FROM CURRENT_DATE)::int,
        EXTRACT(MONTH FROM CURRENT_DATE)::int,
        LEAST(
            rb.billing_day,
            EXTRACT(DAY FROM (DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month - 1 day'))::int
        )
    ),
    'pending'
FROM recurring_bills rb
WHERE rb.user_id = $1 AND rb.active = TRUE
ON CONFLICT DO NOTHING;
```

This guarantees all monthly obligations are present as pending transactions the first time the app is accessed in a given month.

**Note on `billing_day` capping:** `MAKE_DATE` raises an error if the day does not exist in the target month (e.g., day 31 in February). Instead of hardcoding a cap of 28, the query dynamically computes the last valid day of the current month using `DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month - 1 day'` and uses that as the ceiling. This means a bill with `billing_day = 31` is generated on January 31, March 31, April 30, and February 28 (or 29 in a leap year) — always the closest valid date, never artificially early.
