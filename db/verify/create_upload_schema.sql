-- Verify db:create_upload_schema on cockroach

BEGIN;

CREATE PROCEDURE upload.sqitch_verify_create_upload_schema() language plpgsql as $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.schemata WHERE schema_name = 'upload') THEN
        RAISE EXCEPTION 'upload schema does not exist';
    END IF;
END;
$$;

CALL upload.sqitch_verify_create_upload_schema();

ROLLBACK;
