CREATE TABLE transactions (
  id                UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id           UUID          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  recurring_bill_id UUID          REFERENCES recurring_bills(id) ON DELETE SET NULL,
  description       VARCHAR(150)  NOT NULL,
  amount            NUMERIC(12,2) NOT NULL CHECK (amount >= 0),
  category          VARCHAR(50)   NOT NULL DEFAULT 'outros',
  due_date          DATE,
  status            VARCHAR(20)   NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid')),
  paid_at           DATE,
  created_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
