CREATE TABLE "wallets" (
                           "id" bigserial PRIMARY KEY,
                           "owner" varchar NOT NULL,
                           "balance" bigint NOT NULL,
                           "currency" varchar NOT NULL,
                           "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
                           "id" bigserial PRIMARY KEY,
                           "wallet_id" bigint NOT NULL,
                           "amount" bigint NOT NULL,
                           "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "payments" (
                            "id" bigserial PRIMARY KEY,
                            "from_wallet_id" bigint NOT NULL,
                            "to_wallet_id" bigint NOT NULL,
                            "amount" bigint NOT NULL,
                            "status" varchar NOT NULL,
                            "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "wallets" ("owner");

CREATE INDEX ON "entries" ("wallet_id");

CREATE INDEX ON "payments" ("from_wallet_id");

CREATE INDEX ON "payments" ("to_wallet_id");

CREATE INDEX ON "payments" ("from_wallet_id", "to_wallet_id");

COMMENT ON COLUMN "entries"."amount" IS '+credit / -debit';

COMMENT ON COLUMN "payments"."amount" IS 'must be > 0';

COMMENT ON COLUMN "payments"."status" IS 'pending | completed | failed';

ALTER TABLE "entries" ADD FOREIGN KEY ("wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("from_wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("to_wallet_id") REFERENCES "wallets" ("id");
