CREATE TABLE "timeframes" (
  "id" int PRIMARY KEY,
  "code" varchar UNIQUE NOT NULL,
  "description" text
);

CREATE TABLE "instruments" (
  "id" serial PRIMARY KEY,
  "symbol" varchar UNIQUE NOT NULL,
  "name" varchar,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "instrument_timeframes" (
  "instrument_id" int,
  "timeframe_id" int,
  "primary" key(instrument_id,timeframe_id)
);

CREATE TABLE "candles" (
  "id" serial PRIMARY KEY,
  "instrument_id" int,
  "timeframe_id" int,
  "timestamp" timestamptz NOT NULL,
  "open" double,
  "high" double,
  "low" double,
  "close" double,
  "volume" bigint,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE UNIQUE INDEX ON "candles" ("instrument_id", "timeframe_id", "timestamp");

COMMENT ON COLUMN "timeframes"."id" IS 'Enum value (e.g., 12)';

COMMENT ON COLUMN "timeframes"."code" IS 'e.g., ''H1''';

COMMENT ON COLUMN "timeframes"."description" IS 'e.g., ''1 Hour''';

COMMENT ON COLUMN "instruments"."symbol" IS 'e.g., ''GAZP@MISX''';

COMMENT ON COLUMN "candles"."timestamp" IS 'Open time of candle';

ALTER TABLE "instrument_timeframes" ADD FOREIGN KEY ("instrument_id") REFERENCES "instruments" ("id");

ALTER TABLE "instrument_timeframes" ADD FOREIGN KEY ("timeframe_id") REFERENCES "timeframes" ("id");

ALTER TABLE "candles" ADD FOREIGN KEY ("instrument_id") REFERENCES "instruments" ("id");

ALTER TABLE "candles" ADD FOREIGN KEY ("timeframe_id") REFERENCES "timeframes" ("id");
