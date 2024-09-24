CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    metric_type TEXT NOT NULL,
    value DOUBLE PRECISION,
    delta BIGINT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_metrics_name ON metrics (name);