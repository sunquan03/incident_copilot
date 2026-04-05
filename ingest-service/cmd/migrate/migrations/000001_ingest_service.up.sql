
BEGIN;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'outbox_status') THEN
CREATE TYPE outbox_status AS ENUM ('pending', 'published', 'failed');
END IF;
END $$;

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE IF NOT EXISTS alerts (
                                      id                  TEXT PRIMARY KEY,
                                      source_id           TEXT NOT NULL,
                                      source_name         TEXT NOT NULL,
                                      message             TEXT NOT NULL,
                                      labels              JSONB NOT NULL DEFAULT '{}'::jsonb,
                                      created_at          BIGINT NOT NULL,
                                      received_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    dedupe_key          TEXT NOT NULL,
    raw_payload         JSONB NOT NULL DEFAULT '{}'::jsonb,
    metadata            JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_alerts_source_id_not_empty CHECK (btrim(source_id) <> ''),
    CONSTRAINT chk_alerts_source_name_not_empty CHECK (btrim(source_name) <> ''),
    CONSTRAINT chk_alerts_message_not_empty CHECK (btrim(message) <> '')
    );

CREATE UNIQUE INDEX IF NOT EXISTS uq_alerts_dedupe_key
    ON alerts (dedupe_key);

CREATE UNIQUE INDEX IF NOT EXISTS uq_alerts_source_id
    ON alerts (source_id);

CREATE INDEX IF NOT EXISTS idx_alerts_source_name
    ON alerts (source_name);

CREATE INDEX IF NOT EXISTS idx_alerts_created_at
    ON alerts (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_alerts_labels_gin
    ON alerts USING GIN (labels);

CREATE INDEX IF NOT EXISTS idx_alerts_metadata_gin
    ON alerts USING GIN (metadata);

CREATE INDEX IF NOT EXISTS idx_alerts_raw_payload_gin
    ON alerts USING GIN (raw_payload);

DROP TRIGGER IF EXISTS trg_alerts_updated_at ON alerts;
CREATE TRIGGER trg_alerts_updated_at
    BEFORE UPDATE ON alerts
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();


CREATE TABLE IF NOT EXISTS log_docs (
                                        id                  TEXT PRIMARY KEY,
                                        source_id           TEXT NOT NULL,
                                        service_name        TEXT NOT NULL,
                                        title               TEXT NOT NULL,
                                        content             TEXT NOT NULL,
                                        source_type         TEXT NOT NULL,
                                        tags                JSONB NOT NULL DEFAULT '[]'::jsonb,
                                        status              TEXT NOT NULL,
                                        created_at          TIMESTAMPTZ NOT NULL,
                                        received_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    dedupe_key          TEXT NOT NULL,
    raw_payload         JSONB NOT NULL DEFAULT '{}'::jsonb,
    metadata            JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_log_docs_source_id_not_empty CHECK (btrim(source_id) <> ''),
    CONSTRAINT chk_log_docs_service_name_not_empty CHECK (btrim(service_name) <> ''),
    CONSTRAINT chk_log_docs_title_not_empty CHECK (btrim(title) <> ''),
    CONSTRAINT chk_log_docs_content_not_empty CHECK (btrim(content) <> ''),
    CONSTRAINT chk_log_docs_source_type_not_empty CHECK (btrim(source_type) <> ''),
    CONSTRAINT chk_log_docs_status_not_empty CHECK (btrim(status) <> '')
    );

CREATE UNIQUE INDEX IF NOT EXISTS uq_log_docs_dedupe_key
    ON log_docs (dedupe_key);

CREATE UNIQUE INDEX IF NOT EXISTS uq_log_docs_source_id
    ON log_docs (source_id);

CREATE INDEX IF NOT EXISTS idx_log_docs_service_name
    ON log_docs (service_name);

CREATE INDEX IF NOT EXISTS idx_log_docs_source_type
    ON log_docs (source_type);

CREATE INDEX IF NOT EXISTS idx_log_docs_status
    ON log_docs (status);

CREATE INDEX IF NOT EXISTS idx_log_docs_created_at
    ON log_docs (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_log_docs_tags_gin
    ON log_docs USING GIN (tags);

CREATE INDEX IF NOT EXISTS idx_log_docs_metadata_gin
    ON log_docs USING GIN (metadata);

CREATE INDEX IF NOT EXISTS idx_log_docs_raw_payload_gin
    ON log_docs USING GIN (raw_payload);

DROP TRIGGER IF EXISTS trg_log_docs_updated_at ON log_docs;
CREATE TRIGGER trg_log_docs_updated_at
    BEFORE UPDATE ON log_docs
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();


CREATE TABLE IF NOT EXISTS incidents (
                                         id                  TEXT PRIMARY KEY,
                                         source_id           TEXT NOT NULL,
                                         service_name        TEXT NOT NULL,
                                         message             TEXT NOT NULL,
                                         tags                JSONB NOT NULL DEFAULT '[]'::jsonb,
                                         status              TEXT NOT NULL,
                                         created_at          TIMESTAMPTZ NOT NULL,
                                         received_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    dedupe_key          TEXT NOT NULL,
    raw_payload         JSONB NOT NULL DEFAULT '{}'::jsonb,
    metadata            JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_incidents_source_id_not_empty CHECK (btrim(source_id) <> ''),
    CONSTRAINT chk_incidents_service_name_not_empty CHECK (btrim(service_name) <> ''),
    CONSTRAINT chk_incidents_message_not_empty CHECK (btrim(message) <> ''),
    CONSTRAINT chk_incidents_status_not_empty CHECK (btrim(status) <> '')
    );

CREATE UNIQUE INDEX IF NOT EXISTS uq_incidents_dedupe_key
    ON incidents (dedupe_key);

CREATE UNIQUE INDEX IF NOT EXISTS uq_incidents_source_id
    ON incidents (source_id);

CREATE INDEX IF NOT EXISTS idx_incidents_service_name
    ON incidents (service_name);

CREATE INDEX IF NOT EXISTS idx_incidents_status
    ON incidents (status);

CREATE INDEX IF NOT EXISTS idx_incidents_created_at
    ON incidents (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_incidents_tags_gin
    ON incidents USING GIN (tags);

CREATE INDEX IF NOT EXISTS idx_incidents_metadata_gin
    ON incidents USING GIN (metadata);

CREATE INDEX IF NOT EXISTS idx_incidents_raw_payload_gin
    ON incidents USING GIN (raw_payload);

DROP TRIGGER IF EXISTS trg_incidents_updated_at ON incidents;
CREATE TRIGGER trg_incidents_updated_at
    BEFORE UPDATE ON incidents
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();



CREATE TABLE IF NOT EXISTS event_outbox (
                                            id                  BIGSERIAL PRIMARY KEY,
                                            aggregate_type      TEXT NOT NULL,
                                            aggregate_id        TEXT NOT NULL,
                                            event_type          TEXT NOT NULL,
                                            payload             JSONB NOT NULL,
                                            status              outbox_status NOT NULL DEFAULT 'pending',
                                            retry_count         INTEGER NOT NULL DEFAULT 0 CHECK (retry_count >= 0),
    error_message       TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at        TIMESTAMPTZ,
    next_retry_at       TIMESTAMPTZ,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_event_outbox_aggregate_type_not_empty CHECK (btrim(aggregate_type) <> ''),
    CONSTRAINT chk_event_outbox_aggregate_id_not_empty CHECK (btrim(aggregate_id) <> ''),
    CONSTRAINT chk_event_outbox_event_type_not_empty CHECK (btrim(event_type) <> '')
    );

CREATE INDEX IF NOT EXISTS idx_event_outbox_status_created_at
    ON event_outbox (status, created_at);

CREATE INDEX IF NOT EXISTS idx_event_outbox_next_retry_at
    ON event_outbox (next_retry_at)
    WHERE next_retry_at IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_event_outbox_aggregate
    ON event_outbox (aggregate_type, aggregate_id);

CREATE INDEX IF NOT EXISTS idx_event_outbox_payload_gin
    ON event_outbox USING GIN (payload);

DROP TRIGGER IF EXISTS trg_event_outbox_updated_at ON event_outbox;
CREATE TRIGGER trg_event_outbox_updated_at
    BEFORE UPDATE ON event_outbox
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

COMMIT;
