package db

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

func CreateUpload(ctx context.Context, id uuid.UUID, partsCount int, size int, mimeType string) error {
	conn, ok := GetConn(ctx)
	if !ok {
		return errors.New("connection not found in context")
	}
	_, err := conn.Exec(ctx, "CALL upload.create_new_upload($1, $2, $3, $4)", id, partsCount, size, mimeType)
	if err != nil {
		return err
	}
	return nil
}
