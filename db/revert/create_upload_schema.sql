-- Revert db:create_upload_schema from cockroach

BEGIN;

DROP SCHEMA upload CASCADE;

COMMIT;
