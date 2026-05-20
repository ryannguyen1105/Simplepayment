-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2026-05-20T05:10:51.809Z

CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

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

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL
);

CREATE INDEX ON "wallets" ("owner");

CREATE UNIQUE INDEX ON "wallets" ("owner", "currency");

CREATE INDEX ON "entries" ("wallet_id");

CREATE INDEX ON "payments" ("from_wallet_id");

CREATE INDEX ON "payments" ("to_wallet_id");

CREATE INDEX ON "payments" ("from_wallet_id", "to_wallet_id");

COMMENT ON COLUMN "entries"."amount" IS '+credit / -debit';

COMMENT ON COLUMN "payments"."amount" IS 'must be > 0';

COMMENT ON COLUMN "payments"."status" IS 'pending | completed | failed';

ALTER TABLE "wallets" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "entries" ADD FOREIGN KEY ("wallet_id") REFERENCES "wallets" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "payments" ADD FOREIGN KEY ("from_wallet_id") REFERENCES "wallets" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "payments" ADD FOREIGN KEY ("to_wallet_id") REFERENCES "wallets" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username") DEFERRABLE INITIALLY IMMEDIATE;
