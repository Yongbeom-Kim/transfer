-- Deploy db:create_upload_tables to cockroach
-- requires: create_upload_schema

BEGIN;

-- XXX Add DDLs here.

CREATE TYPE upload.upload_status AS ENUM ('pending', 'in_progress', 'completed', 'failed');
CREATE TYPE upload.part_status AS ENUM ('pending', 'uploaded', 'failed');


CREATE TABLE upload.uploads (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status upload.upload_status NOT NULL DEFAULT 'pending',
    -- for object storage
    parts_count INT NOT NULL,
    size INT NOT NULL,
    mime_type TEXT NOT NULL
);

CREATE TABLE upload.parts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    upload_id UUID NOT NULL REFERENCES upload.uploads(id) ON UPDATE CASCADE ON DELETE CASCADE,
    part_number INT NOT NULL,
    status upload.part_status NOT NULL,

    -- for object storage
    object_key TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    uploaded_at TIMESTAMP WITH TIME ZONE,
    byte_offset BIGINT,
    byte_size BIGINT,
    sha256 BYTEA
);

COMMIT;