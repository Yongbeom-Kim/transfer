-- Deploy db:create_upload_table_procedures to cockroach
-- requires: create_upload_tables

BEGIN;

-- Function: Create new upload, return upload id
CREATE PROCEDURE upload.create_new_upload(
    p_id UUID,
    p_parts_count INT,
    p_size INT,
    p_mime_type TEXT
) AS $$
DECLARE
    base_object_key TEXT;
BEGIN
    INSERT INTO upload.uploads (id, parts_count, size, mime_type, status)
    VALUES (p_id, p_parts_count, p_size, p_mime_type, 'pending');

    base_object_key := 'upload-' || p_id;
    FOR part IN 0..(p_parts_count-1) LOOP
        INSERT INTO upload.parts (upload_id, part_number, object_key, status)
        VALUES (p_id, part, base_object_key || '-' || part, 'pending');
    END LOOP;
END
$$ LANGUAGE plpgsql;

-- Procedure: Update upload part status
CREATE OR REPLACE PROCEDURE upload.update_part(
    p_upload_id UUID,
    p_part_number INT,
    p_status upload.part_status,
    p_object_key TEXT,
    p_byte_offset BIGINT,
    p_byte_size BIGINT,
    p_sha256 BYTEA
) AS $$
DECLARE
    part_updated UUID := NULL;
BEGIN
    IF p_upload_id IS NULL THEN
        RAISE EXCEPTION 'Upload ID is required';
    END IF;
    IF p_part_number IS NULL THEN
        RAISE EXCEPTION 'Part number is required';
    END IF;
    IF p_status IS NULL THEN
        RAISE EXCEPTION 'Part status is required';
    END IF;
    -- Wow, this is a multivalued dependency.
    IF p_status = 'uploaded' THEN
        IF p_byte_offset IS NULL OR p_byte_offset = 0 THEN
            RAISE EXCEPTION 'Byte offset is required if part status is uploaded';
        END IF;
        IF p_byte_size IS NULL OR p_byte_size = 0 THEN
            RAISE EXCEPTION 'Byte size is required if part status is uploaded';
        END IF;
        IF p_sha256 IS NULL OR p_sha256 = '' THEN
            RAISE EXCEPTION 'SHA256 is required if part status is uploaded';
        END IF;
    END IF;
    
    UPDATE upload.parts
        SET status = p_status,
            object_key = p_object_key,
            byte_offset = p_byte_offset,
            byte_size = p_byte_size,
            sha256 = p_sha256
        WHERE upload_id = p_upload_id AND part_number = p_part_number
        RETURNING upload_id INTO part_updated;
    IF part_updated IS NULL THEN
        RAISE EXCEPTION 'Part status not found: %, %', p_upload_id, p_part_number;
    END IF;
    IF (SELECT COUNT(*) FROM upload.parts WHERE upload_id = p_upload_id AND status != 'uploaded') = 0 THEN
        UPDATE upload.uploads
        SET status = 'completed'
        WHERE id = p_upload_id;
    ELSE
        UPDATE upload.uploads
        SET status = 'in_progress'
        WHERE id = p_upload_id;
    END IF;
END
$$ LANGUAGE plpgsql;

-- Delete upload
CREATE OR REPLACE PROCEDURE upload.delete_upload(
    p_upload_id UUID
) AS $$
DECLARE
    deleted_count UUID := NULL;
BEGIN
    DELETE FROM upload.uploads WHERE id = p_upload_id RETURNING id INTO deleted_count;
    IF deleted_count IS NULL THEN
        RAISE EXCEPTION 'Upload not found: %', p_upload_id;
    END IF;
    -- Upload parts are deleted by ON DELETE CASCADE
END
$$ LANGUAGE plpgsql;

-- Get keys for upload parts
CREATE FUNCTION upload.get_upload_part_object_keys(
    p_upload_id UUID
) RETURNS SETOF TEXT AS $$
    SELECT object_key FROM upload.parts WHERE upload_id = p_upload_id;
$$ LANGUAGE sql;

COMMIT;
