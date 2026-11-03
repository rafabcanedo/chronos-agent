CREATE TABLE recurring_bills (
  id             UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        UUID          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name           VARCHAR(100)  NOT NULL,
  default_amount NUMERIC(12,2) NOT NULL CHECK (default_amount >= 0),
  billing_day    SMALLINT      NOT NULL CHECK (billing_day BETWEEN 1 AND 31),
  category       VARCHAR(50)   NOT NULL DEFAULT 'outros',
  active         BOOLEAN       NOT NULL DEFAULT TRUE,
  created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
