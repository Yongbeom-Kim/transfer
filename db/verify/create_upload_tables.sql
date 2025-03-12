-- Verify db:create_upload_tables on cockroach

BEGIN;

SELECT id, created_at, status, parts_count, size, mime_type
FROM upload.uploads
WHERE 1=0;

SELECT id, upload_id, part_number, status, created_at, uploaded_at, byte_offset, byte_size, sha256
FROM upload.parts
WHERE 1=0;

ROLLBACK;
