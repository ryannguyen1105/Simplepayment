CREATE TABLE "users" (
                         "id" bigserial PRIMARY KEY,
                         "email" varchar UNIQUE NOT NULL,
                         "username" varchar UNIQUE NOT NULL,
                         "hashed_password" varchar NOT NULL,
                         "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE "wallets" (
                           "id" bigserial PRIMARY KEY,
                           "user_id" bigint UNIQUE NOT NULL REFERENCES users(id),
                           "balance" bigint NOT NULL CHECK (balance >= 0),
                           "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE "payments" (
                            "id" bigserial PRIMARY KEY,
                            "from_user_id" bigint NOT NULL REFERENCES users(id),
                            "to_user_id" bigint NOT NULL REFERENCES users(id),
                            "amount" bigint NOT NULL CHECK (amount >= 0),
                            "status" varchar NOT NULL,
                            "created_at" timestamptz NOT NULL DEFAULT now(),
                            CHECK (from_user_id <> to_user_id)
);

CREATE TABLE "entries" (
                           "id" bigserial PRIMARY KEY,
                           "wallet_id" bigint NOT NULL REFERENCES wallets(id),
                           "payment_id" bigint NOT NULL REFERENCES payments(id),
                           "amount" bigint NOT NULL,
                           "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_email ON users(email);

CREATE INDEX idx_wallets_user_id ON wallets(user_id);

CREATE INDEX idx_payments_from_user_id ON payments(from_user_id);
CREATE INDEX idx_payments_to_user_id ON payments(to_user_id);

CREATE INDEX idx_entries_wallet_id ON entries(wallet_id);
CREATE INDEX idx_entries_payment_id ON entries(payment_id);

COMMENT ON COLUMN "wallets"."balance" IS 'must be >= 0';
COMMENT ON COLUMN "payments"."amount" IS 'must be > 0';
COMMENT ON COLUMN "payments"."status" IS 'pending | completed | failed';
COMMENT ON COLUMN "entries"."amount" IS '+credit / -debit';
