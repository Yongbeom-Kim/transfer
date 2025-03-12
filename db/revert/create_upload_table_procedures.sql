-- Revert db:create_upload_table_procedures from cockroach

BEGIN;

DROP PROCEDURE upload.delete_upload;
DROP PROCEDURE upload.update_part_status;
DROP FUNCTION upload.create_new_upload;
DROP FUNCTION upload.get_upload_part_object_keys;

COMMIT;
