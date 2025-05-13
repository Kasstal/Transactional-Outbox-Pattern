CREATE TABLE "orders" (
                          "id" uuid PRIMARY KEY,
                          "type" varchar(50) NOT NULL,
                          "status" varchar(50) NOT NULL,
                          "city" varchar(100) NOT NULL,
                          "subdivision" varchar(100),
                          "price" decimal(12,2) NOT NULL,
                          "platform" varchar(50) NOT NULL,
                          "general_id" uuid NOT NULL,
                          "order_number" varchar(50),
                          "executor" varchar(100),
                          "created_at" timestamptz NOT NULL DEFAULT (now()),
                          "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "order_items" (
                               "id" serial PRIMARY KEY,
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
                            "type" varchar(20) NOT NULL,
                            "sum" decimal(12,2) NOT NULL,
                            "payed" boolean DEFAULT false,
                            "info" text,
                            "contract_number" varchar(50),
                            "credit_data" jsonb,
                            "external_id" varchar(50),
                            "card_data" jsonb,
                            "created_at" timestamptz NOT NULL DEFAULT (now()),
                            "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "outbox_events" (
                                 "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
                                 "aggregate_type" varchar(50) NOT NULL,
                                 "aggregate_id" uuid NOT NULL,
                                 "event_type" varchar(50) NOT NULL,
                                 "payload" jsonb,
                                 "status" varchar(20) NOT NULL DEFAULT 'pending',
                                 "retry_count" integer DEFAULT 0,
                                 "created_at" timestamptz NOT NULL DEFAULT (now()),
                                 "processed_at" timestamptz
);

CREATE TABLE "history" (
                           "id" serial PRIMARY KEY,
                           "type" varchar(50) NOT NULL,
                           "type_id" integer NOT NULL,
                           "old_value" jsonb,
                           "value" jsonb, --NOT NULL,
                           "date" timestamptz NOT NULL DEFAULT (now()),
                           "user_id" varchar(36) NOT NULL,
                           "order_id" uuid NOT NULL
);

CREATE INDEX ON "orders" ("status");

CREATE INDEX ON "orders" ("platform");

CREATE INDEX ON "orders" ("created_at");

CREATE INDEX ON "orders" ("general_id");

CREATE INDEX ON "order_items" ("order_id");

CREATE INDEX ON "order_items" ("product_id");

CREATE INDEX ON "order_items" ("status");

CREATE INDEX ON "payments" ("order_id");

CREATE INDEX ON "payments" ("type");

CREATE INDEX ON "outbox_events" ("aggregate_type", "aggregate_id");

CREATE INDEX ON "outbox_events" ("status");

CREATE INDEX ON "outbox_events" ("created_at");

CREATE INDEX ON "history" ("order_id");

CREATE INDEX ON "history" ("date");



---ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

--ALTER TABLE "payments" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

--ALTER TABLE "history" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");