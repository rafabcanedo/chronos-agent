CREATE UNIQUE INDEX idx_unique_recurring_per_month
ON transactions (user_id, recurring_bill_id, (DATE_TRUNC('month', due_date::timestamp)))
WHERE recurring_bill_id IS NOT NULL;

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_recurring_bills_updated_at
BEFORE UPDATE ON recurring_bills
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_transactions_updated_at
BEFORE UPDATE ON transactions
FOR EACH ROW EXECUTE FUNCTION set_updated_at();
