-- Deploy db:create_upload_table_procedures to cockroach
-- requires: create_upload_tables

BEGIN;

-- Function: Create new upload, return upload id
CREATE FUNCTION upload.create_new_upload(
    parts_count INT,
    size INT,
    mime_type TEXT
)
RETURNS UUID AS $$
DECLARE
    v_id UUID;
    base_object_key TEXT;
BEGIN
    INSERT INTO upload.uploads (parts_count, size, mime_type, status)
    VALUES (parts_count, size, mime_type, 'pending')
    RETURNING id INTO v_id;

    base_object_key := 'upload-' || v_id;
    FOR part IN 0..(parts_count-1) LOOP
        INSERT INTO upload.parts (upload_id, part_number, object_key, status)
        VALUES (v_id, part, base_object_key || '-' || part, 'pending');
    END LOOP;
    RETURN v_id;
END
$$ LANGUAGE plpgsql;

-- Procedure: Update upload part status
CREATE OR REPLACE PROCEDURE upload.update_part_status(
    upload_id UUID,
    part_number INT,
    new_status upload.part_status
) AS $$
BEGIN
    UPDATE upload.parts
    SET status = new_status
    WHERE upload_id = upload_id AND part_number = part_number;

    IF (SELECT COUNT(*) FROM upload.parts WHERE upload_id = upload_id AND status != 'uploaded') = 0 THEN
        UPDATE upload.uploads
        SET status = 'completed'
        WHERE id = upload_id;
    ELSE
        UPDATE upload.uploads
        SET status = 'in_progress'
        WHERE id = upload_id;
    END IF;
END
$$ LANGUAGE plpgsql;

-- Delete upload
CREATE OR REPLACE PROCEDURE upload.delete_upload(
    upload_id UUID
) AS $$
    DELETE FROM upload.uploads WHERE id = upload_id;
    -- Parts are deleted by ON DELETE CASCADE
$$ LANGUAGE sql;

-- Get keys for upload parts
CREATE FUNCTION upload.get_upload_part_object_keys(
    upload_id UUID
) RETURNS SETOF TEXT AS $$
    SELECT object_key FROM upload.parts WHERE upload_id = upload_id;
$$ LANGUAGE sql;

COMMIT;
