package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
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

type ErrPartStatusNotFound struct {
	UploadID   uuid.UUID
	PartNumber int
}

func (e ErrPartStatusNotFound) Error() string {
	return fmt.Sprintf("part status not found: %s, %d", e.UploadID, e.PartNumber)
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

func UpdateUploadPartStatus(ctx context.Context, uploadID uuid.UUID, partNumber int, status PartStatus) error {
	conn, ok := GetConn(ctx)
	if !ok {
		return errors.New("connection not found in context")
	}
	_, err := conn.Exec(ctx, "CALL upload.update_part_status($1, $2, $3)", uploadID, partNumber, status)
	if err != nil {
		if strings.Contains(err.Error(), "Part status not found:") {
			return ErrPartStatusNotFound{UploadID: uploadID, PartNumber: partNumber}
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

// func GetUpload
// func GetUploadParts
