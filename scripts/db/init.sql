CREATE TABLE IF NOT EXISTS honeypot_logs (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    source_ip INET NOT NULL,
    source_port INTEGER NOT NULL,
    target_port INTEGER NOT NULL,
    protocol VARCHAR(10) NOT NULL,
    raw_payload BYTEA,
    
    CONSTRAINT valid_ports CHECK (source_port >= 0 AND target_port >= 0)
);

CREATE INDEX idx_honeypot_logs_timestamp ON honeypot_logs(timestamp);
CREATE INDEX idx_honeypot_logs_source_ip ON honeypot_logs(source_ip);