DROP INDEX IF EXISTS links2_seed_job_ingestion_timestamp_idx;
DROP INDEX IF EXISTS links2_seq_num_idx;

CREATE INDEX IF EXISTS links2_seed_job_ingestion_ordered_idx ON links2(seed_job_ingest_timestamp DESC NULLS LAST, source_id);
