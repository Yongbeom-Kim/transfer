-- Verify db:create_upload_table_procedures on cockroach

BEGIN;

CREATE FUNCTION upload.sqitch_get_procedure(
    routine_name TEXT,
    routine_type TEXT
) RETURNS TEXT AS $$
BEGIN
    RETURN (SELECT routine_definition FROM information_schema.routines WHERE routine_name = routine_name AND routine_type = routine_type);
END;
$$ LANGUAGE plpgsql;

CREATE PROCEDURE upload.sqitch_verify_create_upload_table_procedures() as $$ 
BEGIN
    IF NOT EXISTS(SELECT upload.sqitch_get_procedure('create_new_upload', 'PROCEDURE')) THEN
        RAISE EXCEPTION 'CREATE_NEW_UPLOAD PROCEDURE DOES NOT EXIST';
    END IF;
    IF NOT EXISTS(SELECT upload.sqitch_get_procedure('update_part', 'PROCEDURE')) THEN
        RAISE EXCEPTION 'UPDATE_PART PROCEDURE DOES NOT EXIST';
    END IF;
    IF NOT EXISTS(SELECT upload.sqitch_get_procedure('delete_upload', 'PROCEDURE')) THEN
        RAISE EXCEPTION 'DELETE_UPLOAD PROCEDURE DOES NOT EXIST';
    END IF;
    IF NOT EXISTS(SELECT upload.sqitch_get_procedure('get_upload_part_object_keys', 'FUNCTION')) THEN
        RAISE EXCEPTION 'GET_UPLOAD_PART_OBJECT_KEYS FUNCTION DOES NOT EXIST';
    END IF;
    RETURN;
END;
$$ language plpgsql;

ROLLBACK;
