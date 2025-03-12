package db

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var ErrUploadAlreadyExists = errors.New("upload already exists")

func CreateUpload(ctx context.Context, id uuid.UUID, partsCount int, size int, mimeType string) error {
	conn, ok := GetConn(ctx)
	if !ok {
		return errors.New("connection not found in context")
	}
	_, err := conn.Exec(ctx, "CALL upload.create_new_upload($1, $2, $3, $4)", id, partsCount, size, mimeType)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"uploads_pkey\"") {
			return ErrUploadAlreadyExists
		}
		return err
	}
	return nil
}
