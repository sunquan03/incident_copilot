BEGIN;

DROP TRIGGER IF EXISTS trg_event_outbox_updated_at ON event_outbox;
DROP TRIGGER IF EXISTS trg_incidents_updated_at ON incidents;
DROP TRIGGER IF EXISTS trg_log_docs_updated_at ON log_docs;
DROP TRIGGER IF EXISTS trg_alerts_updated_at ON alerts;

DROP INDEX IF EXISTS idx_event_outbox_payload_gin;
DROP INDEX IF EXISTS idx_event_outbox_aggregate;
DROP INDEX IF EXISTS idx_event_outbox_next_retry_at;
DROP INDEX IF EXISTS idx_event_outbox_status_created_at;

DROP INDEX IF EXISTS idx_incidents_raw_payload_gin;
DROP INDEX IF EXISTS idx_incidents_metadata_gin;
DROP INDEX IF EXISTS idx_incidents_tags_gin;
DROP INDEX IF EXISTS idx_incidents_created_at;
DROP INDEX IF EXISTS idx_incidents_status;
DROP INDEX IF EXISTS idx_incidents_service_name;
DROP INDEX IF EXISTS uq_incidents_source_id;
DROP INDEX IF EXISTS uq_incidents_dedupe_key;

DROP INDEX IF EXISTS idx_log_docs_raw_payload_gin;
DROP INDEX IF EXISTS idx_log_docs_metadata_gin;
DROP INDEX IF EXISTS idx_log_docs_tags_gin;
DROP INDEX IF EXISTS idx_log_docs_created_at;
DROP INDEX IF EXISTS idx_log_docs_status;
DROP INDEX IF EXISTS idx_log_docs_source_type;
DROP INDEX IF EXISTS idx_log_docs_service_name;
DROP INDEX IF EXISTS uq_log_docs_source_id;
DROP INDEX IF EXISTS uq_log_docs_dedupe_key;

-- Drop indexes on alerts
DROP INDEX IF EXISTS idx_alerts_raw_payload_gin;
DROP INDEX IF EXISTS idx_alerts_metadata_gin;
DROP INDEX IF EXISTS idx_alerts_labels_gin;
DROP INDEX IF EXISTS idx_alerts_created_at;
DROP INDEX IF EXISTS idx_alerts_source_name;
DROP INDEX IF EXISTS uq_alerts_source_id;
DROP INDEX IF EXISTS uq_alerts_dedupe_key;

DROP TABLE IF EXISTS event_outbox;
DROP TABLE IF EXISTS incidents;
DROP TABLE IF EXISTS log_docs;
DROP TABLE IF EXISTS alerts;

DROP FUNCTION IF EXISTS set_updated_at();

DROP TYPE IF EXISTS outbox_status;

COMMIT;