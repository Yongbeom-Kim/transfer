-- Deploy db:create_upload_schema to cockroach

BEGIN;

CREATE SCHEMA upload;

COMMIT;
