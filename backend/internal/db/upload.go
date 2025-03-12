package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Upload struct {
	ID         uuid.UUID
	PartsCount int
	Size       int
	MimeType   string
	Status     UploadStatus
	CreatedAt  time.Time
}

type Part struct {
	ID         uuid.UUID
	UploadID   uuid.UUID
	PartNumber int
	Status     PartStatus
	ObjectKey  string
	CreatedAt  time.Time
	UploadedAt *time.Time
	ByteOffset *int64
	ByteSize   *int64
	Sha256     *[]byte
}

type UploadStatus string

const (
	UploadStatusPending    UploadStatus = "pending"
	UploadStatusInProgress UploadStatus = "in_progress"
	UploadStatusCompleted  UploadStatus = "completed"
	UploadStatusFailed     UploadStatus = "failed"
)

type PartStatus string

const (
	PartStatusPending  PartStatus = "pending"
	PartStatusUploaded PartStatus = "uploaded"
	PartStatusFailed   PartStatus = "failed"
)

type ErrUploadAlreadyExists struct {
	UploadID uuid.UUID
}

func (e ErrUploadAlreadyExists) Error() string {
	return fmt.Sprintf("upload already exists: %s", e.UploadID)
}

type ErrPartNotFound struct {
	UploadID   uuid.UUID
	PartNumber int
}

func (e ErrPartNotFound) Error() string {
	return fmt.Sprintf("part not found: %s, %d", e.UploadID, e.PartNumber)
}

type ErrUploadNotFound struct {
	UploadID uuid.UUID
}

func (e ErrUploadNotFound) Error() string {
	return fmt.Sprintf("upload not found: %s", e.UploadID)
}

func CreateUpload(ctx context.Context, id uuid.UUID, partsCount int, size int, mimeType string) error {
	conn, ok := GetConn(ctx)
	if !ok {
		return errors.New("connection not found in context")
	}
	_, err := conn.Exec(ctx, "CALL upload.create_new_upload($1, $2, $3, $4)", id, partsCount, size, mimeType)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"uploads_pkey\"") {
			return ErrUploadAlreadyExists{UploadID: id}
		}
		return err
	}
	return nil
}

func UpdateUploadPart(ctx context.Context, newPart Part) error {
	conn, ok := GetConn(ctx)
	if !ok {
		return errors.New("connection not found in context")
	}
	_, err := conn.Exec(ctx, "CALL upload.update_part($1, $2, $3, $4, $5, $6, $7)", newPart.UploadID, newPart.PartNumber, newPart.Status, newPart.ObjectKey, newPart.ByteOffset, newPart.ByteSize, newPart.Sha256)
	if err != nil {
		if strings.Contains(err.Error(), "Part status not found:") {
			return ErrPartNotFound{UploadID: newPart.UploadID, PartNumber: newPart.PartNumber}
		}
		return err
	}
	return nil
}

func DeleteUpload(ctx context.Context, uploadID uuid.UUID) error {
	conn, ok := GetConn(ctx)
	if !ok {
		return errors.New("connection not found in context")
	}
	_, err := conn.Exec(ctx, "CALL upload.delete_upload($1)", uploadID)
	if err != nil {
		if strings.Contains(err.Error(), "Upload not found:") {
			return ErrUploadNotFound{UploadID: uploadID}
		}
		return err
	}
	return nil
}

func GetUpload(ctx context.Context, uploadID uuid.UUID) (*Upload, error) {
	conn, ok := GetConn(ctx)
	if !ok {
		return nil, errors.New("connection not found in context")
	}
	var upload Upload
	row := conn.QueryRow(ctx, "SELECT id, created_at, status, parts_count, size, mime_type FROM upload.uploads WHERE id = $1", uploadID)
	err := row.Scan(&upload.ID, &upload.CreatedAt, &upload.Status, &upload.PartsCount, &upload.Size, &upload.MimeType)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, ErrUploadNotFound{UploadID: uploadID}
		}
		return nil, err
	}
	return &upload, nil
}

func GetUploadParts(ctx context.Context, uploadID uuid.UUID) ([]Part, error) {
	conn, ok := GetConn(ctx)
	if !ok {
		return nil, errors.New("connection not found in context")
	}
	rows, err := conn.Query(ctx, "SELECT id, upload_id, part_number, status, object_key, created_at, uploaded_at, byte_offset, byte_size, sha256 FROM upload.parts WHERE upload_id = $1", uploadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	parts := []Part{}
	for rows.Next() {
		var part Part
		err := rows.Scan(&part.ID, &part.UploadID, &part.PartNumber, &part.Status, &part.ObjectKey, &part.CreatedAt, &part.UploadedAt, &part.ByteOffset, &part.ByteSize, &part.Sha256)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}
	if len(parts) == 0 {
		return nil, ErrUploadNotFound{UploadID: uploadID}
	}
	return parts, nil
}
