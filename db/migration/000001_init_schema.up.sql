CREATE TYPE "payment_type" AS ENUM (
  'cash_at_shop',
  'cash_to_courier',
  'card',
  'card_online',
  'credit',
  'bonuses',
  'cashless',
  'prepayment'
);

CREATE TABLE "orders" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "type" varchar(50) NOT NULL,
  "status" varchar(50) NOT NULL,
  "city" varchar(100) NOT NULL,
  "subdivision" varchar(100),
  "price" decimal(12,2) NOT NULL,
  "platform" varchar(50) NOT NULL,
  "general_id" uuid NOT NULL,
  "order_number" varchar(50) UNIQUE NOT NULL,
  "executor" varchar(100),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "history" (
  "id" serial PRIMARY KEY,
  "type" varchar(50) NOT NULL,
  "type_id" integer NOT NULL,
  "old_value" bytea,
  "value" bytea NOT NULL,
  "date" timestamptz NOT NULL DEFAULT (now()),
  "user_id" varchar(36) NOT NULL,
  "order_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "order_items" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "product_id" varchar(36) NOT NULL,
  "external_id" varchar(50),
  "status" varchar(50) NOT NULL,
  "base_price" decimal(12,2) NOT NULL,
  "price" decimal(12,2) NOT NULL,
  "earned_bonuses" decimal(12,2) DEFAULT 0,
  "spent_bonuses" decimal(12,2) DEFAULT 0,
  "gift" boolean DEFAULT false,
  "owner_id" varchar(36),
  "delivery_id" varchar(36),
  "shop_assistant" varchar(100),
  "warehouse" varchar(100),
  "order_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "payments" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "order_id" uuid NOT NULL,
  "type" payment_type NOT NULL,
  "sum" decimal(12,2) NOT NULL,
  "payed" boolean DEFAULT false,
  "info" text,
  "contract_number" varchar(50),
  "external_id" varchar(50),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "credit_data" (
  "id" serial PRIMARY KEY,
  "payment_id" uuid NOT NULL,
  "bank" varchar(100) NOT NULL,
  "type" varchar(50) NOT NULL,
  "number_of_months" smallint NOT NULL,
  "pay_sum_per_month" decimal(12,2) NOT NULL,
  "broker_id" integer,
  "iin" varchar(12),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "card_payment_data" (
  "id" serial PRIMARY KEY,
  "payment_id" uuid NOT NULL,
  "provider" varchar(100) NOT NULL,
  "transaction_id" varchar(100) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- Indexes and foreign keys remain the same as your original
CREATE INDEX ON "orders" ("status");
CREATE INDEX ON "orders" ("platform");
CREATE INDEX ON "orders" ("created_at");
CREATE INDEX ON "orders" ("general_id");
CREATE INDEX ON "history" ("order_id");
CREATE INDEX ON "history" ("date");
CREATE INDEX ON "order_items" ("order_id");
CREATE INDEX ON "order_items" ("product_id");
CREATE INDEX ON "order_items" ("status");
CREATE INDEX ON "payments" ("order_id");
CREATE INDEX ON "payments" ("type");

ALTER TABLE "history" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
ALTER TABLE "payments" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");
ALTER TABLE "credit_data" ADD FOREIGN KEY ("payment_id") REFERENCES "payments" ("id") ON DELETE CASCADE;
ALTER TABLE "card_payment_data" ADD FOREIGN KEY ("payment_id") REFERENCES "payments" ("id") ON DELETE CASCADE;