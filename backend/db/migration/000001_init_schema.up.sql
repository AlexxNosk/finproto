-- Table of all available timeframes
CREATE TABLE IF NOT EXISTS timeframes (
    id          INTEGER PRIMARY KEY,           -- Enum value (e.g., 12)
    code        VARCHAR NOT NULL UNIQUE,       -- e.g., "H1"
    description TEXT                           -- Optional description
);

-- Pre-fill timeframes (only if not already present)
INSERT INTO timeframes (id, code, description) VALUES
(0,  'UNSPECIFIED', 'Timeframe not specified'),
(1,  'M1',          '1 Minute'),
(5,  'M5',          '5 Minutes'),
(9,  'M15',         '15 Minutes'),
(11, 'M30',         '30 Minutes'),
(12, 'H1',          '1 Hour'),
(13, 'H2',          '2 Hours'),
(15, 'H4',          '4 Hours'),
(17, 'H8',          '8 Hours'),
(19, 'D',           'Day'),
(20, 'W',           'Week'),
(21, 'MN',          'Month'),
(22, 'QR',          'Quarter')
ON CONFLICT (id) DO NOTHING;

-- assets table
CREATE TABLE IF NOT EXISTS finam_assets (
    id           SERIAL PRIMARY KEY,
    ticker       VARCHAR NOT NULL,
    symbol       VARCHAR NOT NULL UNIQUE,        -- e.g., "GAZP@MISX"
    name         VARCHAR,
    mic          VARCHAR NOT NULL,               -- Exchange code
    type         VARCHAR NOT NULL,               -- asset type
    external_id  VARCHAR NOT NULL UNIQUE,        -- Original ID from external source
    updated_at   TIMESTAMPTZ DEFAULT now()
);

-- Track per-asset+timeframe table creation
CREATE TABLE IF NOT EXISTS asset_tables (
    id             SERIAL PRIMARY KEY,
    asset_id  INT NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    timeframe_id   INT NOT NULL REFERENCES timeframes(id),
    table_name     TEXT NOT NULL UNIQUE,
    created_at     TIMESTAMPTZ DEFAULT now(),
    updated_at     TIMESTAMPTZ DEFAULT now(),
    UNIQUE (asset_id, timeframe_id)
);

-- Track per-asset+timeframe data tables creation
CREATE TABLE IF NOT EXISTS data_tables (
    id                      SERIAL PRIMARY KEY,
    asset_id           INT NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    timeframe_id            INT NOT NULL REFERENCES timeframes(id),
    asset_table_id     INT NOT NULL REFERENCES asset_tables(id),
    table_name              TEXT NOT NULL UNIQUE,
    purpose                 TEXT,                  
    created_at              TIMESTAMPTZ DEFAULT now(),
    updated_at              TIMESTAMPTZ DEFAULT now(),

    UNIQUE (asset_id, timeframe_id, table_name)
);


-- Historical bars table
-- CREATE TABLE bars (
--     id            SERIAL PRIMARY KEY,
    -- asset_id INT NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    -- timeframe_id  INT NOT NULL REFERENCES timeframes(id),
    -- timestamp     TIMESTAMPTZ NOT NULL,            
    -- open          DOUBLE PRECISION,                -- Nullable
    -- high          DOUBLE PRECISION,
    -- low           DOUBLE PRECISION,
    -- close         DOUBLE PRECISION,
    -- volume        BIGINT --,
    -- created_at    TIMESTAMPTZ DEFAULT now(),
    -- updated_at    TIMESTAMPTZ,
    -- UNIQUE (asset_id, timeframe_id, timestamp)
-- );


