CREATE TABLE IF NOT EXISTS accounts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    VARCHAR(100) NOT NULL,
    currency    VARCHAR(3)   NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key VARCHAR(100)    NOT NULL UNIQUE,
    from_account_id UUID            NOT NULL REFERENCES accounts(id),
    to_account_id   UUID            NOT NULL REFERENCES accounts(id),
    amount          NUMERIC(20,8)   NOT NULL CHECK (amount > 0),
    currency        VARCHAR(3)      NOT NULL,
    status          VARCHAR(20)     NOT NULL DEFAULT 'completed',
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS entries (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id     UUID          NOT NULL REFERENCES accounts(id),
    transaction_id UUID          NOT NULL,
    amount         NUMERIC(20,8) NOT NULL CHECK (amount > 0),
    direction      VARCHAR(6)    NOT NULL CHECK (direction IN ('debit','credit')),
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_entries_account_id      ON entries(account_id);
CREATE INDEX idx_entries_transaction_id  ON entries(transaction_id);
CREATE INDEX idx_transactions_idempotency_key ON transactions(idempotency_key);
CREATE INDEX idx_transactions_from_account    ON transactions(from_account_id);
CREATE INDEX idx_transactions_to_account      ON transactions(to_account_id);
