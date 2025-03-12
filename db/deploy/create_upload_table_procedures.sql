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
CREATE OR REPLACE PROCEDURE upload.update_part_status(
    p_upload_id UUID,
    p_part_number INT,
    p_new_status upload.part_status
) AS $$
BEGIN
    UPDATE upload.parts
    SET status = p_new_status
    WHERE upload_id = p_upload_id AND part_number = p_part_number;

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
    DELETE FROM upload.uploads WHERE id = p_upload_id;
    -- Parts are deleted by ON DELETE CASCADE
$$ LANGUAGE sql;

-- Get keys for upload parts
CREATE FUNCTION upload.get_upload_part_object_keys(
    p_upload_id UUID
) RETURNS SETOF TEXT AS $$
    SELECT object_key FROM upload.parts WHERE upload_id = p_upload_id;
$$ LANGUAGE sql;

COMMIT;
