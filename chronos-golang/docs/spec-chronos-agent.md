# Personal Finance App — Backend & Supabase Technical Specification

## Overview
This document specifies the backend system and database architecture for a lightweight personal finance tracking application. The backend is written in **Go** and exposes a RESTful JSON API. Data persistence is handled by **PostgreSQL hosted on Supabase**. Natural language interactions are parsed using **Abacus RouteLLM** via an OpenAI-compatible HTTP interface.

## 1. Technology Stack & Infrastructure

- **Language:** Go (modern standard toolchain)
- **Database:** PostgreSQL (Supabase Free Tier)
- **Database Driver:** `pgx` (recommended) or standard `database/sql` + `lib/pq`
- **Authentication:** Stateless JWT (Bearer Token)
- **Password Hashing:** `bcrypt`
- **AI Integration:** Abacus RouteLLM API (`https://routellm.abacus.ai/v1/chat/completions`)
- **Hosting Target:** Containerized service on free-tier cloud host (e.g., Render, Fly.io, Koyeb)
- **Obsidian Toolkit:** We need to use our toolkit(`https://github.com/rafabcanedo/obsidian-golang-toolkit`)

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
password_hash
TEXT
NOT NULL
bcrypt hash of the user password
created_at
TIMESTAMPTZ
NOT NULL DEFAULT NOW()
Record creation timestamp
Note: User registration is not exposed via an open public endpoint. Users are created via seed scripts or directly through the Supabase SQL editor.

### 2.2 Table: recurring_bills
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
due_day
SMALLINT
NOT NULL CHECK (due_day BETWEEN 1 AND 31)
Day of the month the bill is due
category
VARCHAR(50)
NOT NULL DEFAULT 'outros'
Expense category
active
BOOLEAN
NOT NULL DEFAULT TRUE
Whether to generate instances each month
created_at
TIMESTAMPTZ
NOT NULL DEFAULT NOW()
Record creation timestamp

### 2.3 Table: transactions
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
Category
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


### 2.4 Integrity & Deduplication Rules
Monthly recurring deduplication: A unique compound index prevents generating duplicate transactions for the same recurring bill within the same calendar month:
sql
Copy
CREATE UNIQUE INDEX idx_unique_recurring_per_month 
ON transactions (user_id, recurring_bill_id, (DATE_TRUNC('month', due_date)))
WHERE recurring_bill_id IS NOT NULL;

Monetary values: All currencies strictly use NUMERIC(12,2) to prevent floating-point rounding inaccuracies.
Overdue state: There is no dedicated overdue database column. A record is computed as overdue by the application when status = 'pending' AND due_date < CURRENT_DATE.

## 3. Go Project Structure & Architectural Layers
The service follows a domain-driven, layered architecture without unnecessary enterprise abstractions.

text
Copy
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

Layer Responsibilities
Handlers (HTTP): Decode JSON payloads, extract user context (user_id from JWT), enforce basic validation, return JSON responses.
Services (Business Logic): Contain domain operations (e.g., fuzzy-matching a bill description, synchronizing status = 'paid' with paid_at, generating monthly bills).
Repositories (Data Access): Execute raw or prepared SQL queries against Supabase. Enforce WHERE user_id = $1 on every query.
Chat Orchestrator (chat package): Calls the llm package to get a structured JSON action, performs validation, and delegates execution to the appropriate service.

## 4. API Endpoints Specification
All endpoints except POST /api/auth/login require the Authorization: Bearer <JWT_TOKEN> header.

### 4.1 Authentication
POST /api/auth/login
Body: {"username": "string", "password": "string"}
Returns: 200 OK with {"token": "jwt_string_here"} or 401 Unauthorized.
Token Expiry: 30 days. Contains user_id in claims.

### 4.2 Transactions
GET /api/transactions?month=YYYY-MM&status=pending|paid
Fetches user transactions filtered by calendar month and/or payment state.
PATCH /api/transactions/:id
Updates mutable fields (e.g., marking as paid: {"status": "paid", "paid_at": "YYYY-MM-DD"}).
DELETE /api/transactions/:id
Permanently removes a transaction record.

### 4.3 Recurring Bills (Templates)
GET /api/recurring
Lists all active recurring bill templates.
POST /api/recurring
Creates a new recurring template (name, default_amount, due_day, category).
PUT /api/recurring/:id
Updates an existing recurring template (changes apply only to future generated transactions).

### 4.4 Chat & NLP Processing
POST /api/chat
Body: {"message": "string"}
Flow:
Extracts user_id from JWT context.
Sends prompt and current server date (YYYY-MM-DD) to Abacus RouteLLM.
Parses structured JSON response from the LLM.
Executes the corresponding database operation.
Returns confirmation message or a descriptive error to the client.

## 5. Abacus RouteLLM Integration & Prompt Contract
The backend uses standard HTTP POST to communicate with Abacus RouteLLM.

Endpoint: https://routellm.abacus.ai/v1/chat/completions
Model: route-llm
Headers: Authorization: Bearer <ABACUS_API_KEY>, Content-Type: application/json
Prompting Strategy
The backend sets a strict system prompt instructing the model to output only raw JSON with no markdown fences or prose commentary:

text
Copy
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
  "due_day": integer (1-31) or null,
  "category": string or null
}

Action Handling Logic in Backend
action == "pay":
Go runs an ILIKE / pattern match for pending transactions in the current month matching description:
sql
Copy
SELECT id FROM transactions 
WHERE user_id = $1 
  AND status = 'pending' 
  AND description ILIKE '%' || $2 || '%' 
  AND DATE_TRUNC('month', due_date) = DATE_TRUNC('month', $3::date)
LIMIT 1;

If found: Updates status = 'paid', paid_at = date, returns success.
If not found: Does NOT auto-create a bill. Returns HTTP 404 or error message: "Conta não encontrada para este mês."
action == "create_expense":
Inserts immediately into transactions with status = 'paid', paid_at = date, recurring_bill_id = NULL. Returns success.
action == "create_recurring":
Inserts into recurring_bills (name, default_amount, due_day, category). Returns success.
action == "query":
Backend performs aggregate query (SUM(amount) grouped by category or status) for the requested month, passes numbers back to the LLM or formats directly in Go, and returns the answer.

## 6. Monthly Recurring Bills Generation (Backend Logic)
To ensure bills for the new month exist without needing external cron jobs, the backend uses lazy evaluation:

When a user requests GET /api/transactions or POST /api/chat, the backend checks if bills for the current month (DATE_TRUNC('month', NOW())) have already been populated for this user_id.
If unpopulated, it executes an idempotent batch INSERT:
sql
Copy
INSERT INTO transactions (user_id, recurring_bill_id, description, amount, category, due_date, status)
SELECT 
    rb.user_id,
    rb.id,
    rb.name,
    rb.default_amount,
    rb.category,
    MAKE_DATE(EXTRACT(YEAR FROM CURRENT_DATE)::int, EXTRACT(MONTH FROM CURRENT_DATE)::int, LEAST(rb.due_day, 28)),
    'pending'
FROM recurring_bills rb
WHERE rb.user_id = $1 AND rb.active = TRUE
ON CONFLICT DO NOTHING;

This guarantees all monthly obligations are present as pending transactions the first time the app is accessed in a given month.
