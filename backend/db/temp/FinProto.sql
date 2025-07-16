-- Table of all available timeframes (based on your enum)
CREATE TABLE timeframes (
    id          INTEGER PRIMARY KEY,          -- Enum value (e.g., 12)
    code        VARCHAR NOT NULL UNIQUE,      -- e.g., "H1"
    description TEXT                          -- Optional description (e.g., "1 Hour")
);

-- Pre-fill the timeframes table
INSERT INTO timeframes (id, code, description) VALUES
(0,  'UNSPECIFIED', 'Timeframe not specified'),
(1,  'M1',  '1 Minute'),
(5,  'M5',  '5 Minutes'),
(9,  'M15', '15 Minutes'),
(11, 'M30', '30 Minutes'),
(12, 'H1',  '1 Hour'),
(13, 'H2',  '2 Hours'),
(15, 'H4',  '4 Hours'),
(17, 'H8',  '8 Hours'),
(19, 'D',   'Day'),
(20, 'W',   'Week'),
(21, 'MN',  'Month'),
(22, 'QR',  'Quarter');

-- Instruments table (e.g., GAZP@MISX)
CREATE TABLE instruments (
    id          SERIAL PRIMARY KEY,
    symbol      VARCHAR NOT NULL UNIQUE,      -- e.g., "GAZP@MISX"
    name        VARCHAR,                      -- Optional descriptive name
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ
);

-- Many-to-many relation: which timeframes each instrument supports
CREATE TABLE instrument_timeframes (
    instrument_id INT NOT NULL REFERENCES instruments(id) ON DELETE CASCADE,
    timeframe_id  INT NOT NULL REFERENCES timeframes(id),
    PRIMARY KEY (instrument_id, timeframe_id)
);

-- Historical candles table
CREATE TABLE candles (
    id            SERIAL PRIMARY KEY,
    instrument_id INT NOT NULL REFERENCES instruments(id) ON DELETE CASCADE,
    timeframe_id  INT NOT NULL REFERENCES timeframes(id),
    timestamp     TIMESTAMPTZ NOT NULL,            -- Open time of the candle
    open          DOUBLE PRECISION,                -- Nullable
    high          DOUBLE PRECISION,
    low           DOUBLE PRECISION,
    close         DOUBLE PRECISION,
    volume        BIGINT,
    created_at    TIMESTAMPTZ DEFAULT now(),
    updated_at    TIMESTAMPTZ,
    UNIQUE (instrument_id, timeframe_id, timestamp)
);
