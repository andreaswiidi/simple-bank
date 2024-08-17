CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "fullname" varchar NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "transaction_history" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "ammount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers_history" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "ammount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "transaction_history" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers_history" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers_history" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

CREATE INDEX ON "accounts" ("username");

CREATE INDEX ON "transaction_history" ("account_id");

CREATE INDEX ON "transfers_history" ("from_account_id");

CREATE INDEX ON "transfers_history" ("to_account_id");

CREATE INDEX ON "transfers_history" ("from_account_id", "to_account_id");