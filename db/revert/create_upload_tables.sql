-- Revert db:create_upload_tables from cockroach

BEGIN;

DROP TABLE upload.parts;
DROP TABLE upload.uploads;
DROP TYPE upload.part_status;
DROP TYPE upload.upload_status;

COMMIT;
