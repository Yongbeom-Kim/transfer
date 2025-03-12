package db

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func SetupTest(t *testing.T) (context.Context, *pgx.Tx, func()) {
	dbpool, cleanup, err := InitDBPool()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	tx, err := dbpool.Begin(context.Background())
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	ctx := WithTx(context.Background(), &tx)
	return ctx, &tx, func() {
		tx.Rollback(context.Background())
		cleanup()
	}
}

func TestCreateUpload(t *testing.T) {
	// Initialize the database connection
	id := uuid.New()
	ctx, tx, cleanup := SetupTest(t)
	defer cleanup()

	err := CreateUpload(ctx, id, 1, 1024, "image/jpeg")
	if err != nil {
		t.Fatalf("Failed to create upload: %v", err)
	}

	// Check if the upload exists in the database
	var upload struct {
		ID         string
		CreatedAt  time.Time
		Status     string
		PartsCount int
		Size       int
		MimeType   string
	}
	err = (*tx).QueryRow(context.Background(), "(SELECT * FROM upload.uploads WHERE id = $1)", id).Scan(&upload.ID, &upload.CreatedAt, &upload.Status, &upload.PartsCount, &upload.Size, &upload.MimeType)
	if err != nil {
		t.Fatalf("Failed to check if upload exists: %v", err)
	}
	if upload.ID != id.String() {
		t.Fatalf("Upload ID does not match")
	}
	if upload.Status != "pending" {
		t.Fatalf("Upload status is not pending")
	}
	if upload.PartsCount != 1 {
		t.Fatalf("Expected 1 upload part, got %d", upload.PartsCount)
	}
	if upload.Size != 1024 {
		t.Fatalf("Expected size to be 1024, got %d", upload.Size)
	}
	if upload.MimeType != "image/jpeg" {
		t.Fatalf("Expected mime type to be image/jpeg, got %s", upload.MimeType)
	}

	// Check if the upload parts exist in the database
	var partsCount int
	err = (*tx).QueryRow(context.Background(), "SELECT COUNT(*) FROM upload.parts WHERE upload_id = $1", id).Scan(&partsCount)
	if err != nil {
		t.Fatalf("Failed to check if upload parts exist: %v", err)
	}
	if partsCount != 1 {
		t.Fatalf("Expected 1 upload part, got %d", partsCount)
	}

}

func TestCreateUpload_Duplicate(t *testing.T) {
	id := uuid.New()
	ctx, _, cleanup := SetupTest(t)
	defer cleanup()

	err := CreateUpload(ctx, id, 1, 1024, "image/jpeg")
	if err != nil {
		t.Fatalf("Failed to create upload: %v", err)
	}

	// Attempt to create the same upload again
	err = CreateUpload(ctx, id, 1, 1024, "image/jpeg")
	var target ErrUploadAlreadyExists
	if !errors.As(err, &target) {
		t.Fatalf("Expected ErrUploadAlreadyExists, got %v", err)
	}
}

func TestUpdateUploadPart(t *testing.T) {
	id := uuid.New()
	ctx, tx, cleanup := SetupTest(t)
	defer cleanup()

	// Call CreateUpload() with 2 parts
	err := CreateUpload(ctx, id, 2, 2048, "application/pdf")
	if err != nil {
		t.Fatalf("Failed to create upload: %v", err)
	}

	// Verify overall upload is pending
	var uploadStatus string
	err = (*tx).QueryRow(context.Background(), "SELECT status FROM upload.uploads WHERE id = $1", id).Scan(&uploadStatus)
	if err != nil {
		t.Fatalf("Failed to check upload status: %v", err)
	}
	if uploadStatus != "pending" {
		t.Fatalf("Expected upload status to be pending, got %s", uploadStatus)
	}

	byteOffset := int64(1024)
	byteSize := int64(1024)
	sha256 := []byte("1234567890")

	// Change 1 part to uploaded
	err = UpdateUploadPart(ctx, Part{
		UploadID:   id,
		PartNumber: 0,
		Status:     PartStatusUploaded,
		ByteOffset: &byteOffset,
		ByteSize:   &byteSize,
		Sha256:     &sha256,
	})
	if err != nil {
		t.Fatalf("Failed to update part 0 status: %v", err)
	}
	// Verify overall upload is in progress
	err = (*tx).QueryRow(context.Background(), "SELECT status FROM upload.uploads WHERE id = $1", id).Scan(&uploadStatus)
	if err != nil {
		t.Fatalf("Failed to check upload status: %v", err)
	}
	if uploadStatus != "in_progress" {
		t.Fatalf("Expected upload status to be in_progress, got %s", uploadStatus)
	}

	byteOffset2 := int64(2048)
	byteSize2 := int64(2048)
	sha2562 := []byte("1234567890")

	// Change 2nd part to uploaded
	err = UpdateUploadPart(ctx, Part{
		UploadID:   id,
		PartNumber: 1,
		Status:     PartStatusUploaded,
		ByteOffset: &byteOffset2,
		ByteSize:   &byteSize2,
		Sha256:     &sha2562,
	})
	if err != nil {
		t.Fatalf("Failed to update part 1 status: %v", err)
	}
	// Verify both parts are uploaded
	var partStatus string
	for partNumber := 0; partNumber < 2; partNumber++ {
		err = (*tx).QueryRow(context.Background(), "SELECT status FROM upload.parts WHERE upload_id = $1 AND part_number = $2", id, partNumber).Scan(&partStatus)
		if err != nil {
			t.Fatalf("Failed to check part %d status: %v", partNumber, err)
		}
		if partStatus != "uploaded" {
			t.Fatalf("Expected part %d status to be uploaded, got %s", partNumber, partStatus)
		}
	}

	// Verify overall upload is completed
	err = (*tx).QueryRow(context.Background(), "SELECT status FROM upload.uploads WHERE id = $1", id).Scan(&uploadStatus)
	if err != nil {
		t.Fatalf("Failed to check upload status: %v", err)
	}
	if uploadStatus != "completed" {
		t.Fatalf("Expected upload status to be completed, got %s", uploadStatus)
	}
}

func TestUpdateUploadPart_UploadedNoInfo(t *testing.T) {
	id := uuid.New()
	ctx, _, cleanup := SetupTest(t)
	defer cleanup()
	// Create upload with 1 part
	err := CreateUpload(ctx, id, 1, 1024, "image/jpeg")
	if err != nil {
		t.Fatalf("Failed to create upload: %v", err)
	}

	// Attempt to update the uploaded part with nil byte offset, byte size, and sha256
	err = UpdateUploadPart(ctx, Part{
		UploadID:   id,
		PartNumber: 0,
		Status:     PartStatusUploaded,
		ByteOffset: nil,
		ByteSize:   nil,
		Sha256:     nil,
	})
	if err == nil {
		t.Fatalf("Expected error when updating part with nil byte offset, byte size, and sha256, but got none")
	}

}

func TestUpdateUploadPartStatus_NoStatusToUpdate(t *testing.T) {
	id := uuid.New()
	ctx, _, cleanup := SetupTest(t)
	defer cleanup()
	byteOffset := int64(1024)
	byteSize := int64(1024)
	sha256 := []byte("1234567890")
	// Attempt to update a non-existent part status
	err := UpdateUploadPart(ctx, Part{
		UploadID:   id,
		PartNumber: 3,
		Status:     PartStatusUploaded,
		ByteOffset: &byteOffset,
		ByteSize:   &byteSize,
		Sha256:     &sha256,
	})
	var target ErrPartNotFound
	if !errors.As(err, &target) {
		t.Fatalf("Expected ErrPartStatusNotFound, got %v", err)
	}
}

func TestDeleteUpload(t *testing.T) {
	id := uuid.New()
	ctx, tx, cleanup := SetupTest(t)
	defer cleanup()

	// Create upload
	err := CreateUpload(ctx, id, 2, 1024, "application/octet-stream")
	if err != nil {
		t.Fatalf("Failed to create upload: %v", err)
	}

	// Delete upload
	err = DeleteUpload(ctx, id)
	if err != nil {
		t.Fatalf("Failed to delete upload: %v", err)
	}

	// Ensure upload is not present
	var uploadCount int
	err = (*tx).QueryRow(context.Background(), "SELECT COUNT(*) FROM upload.uploads WHERE id = $1", id).Scan(&uploadCount)
	if err != nil {
		t.Fatalf("Failed to check upload presence: %v", err)
	}
	if uploadCount != 0 {
		t.Fatalf("Expected upload to be deleted, but it is still present")
	}

	// Ensure upload part is not present
	var partCount int
	err = (*tx).QueryRow(context.Background(), "SELECT COUNT(*) FROM upload.parts WHERE upload_id = $1", id).Scan(&partCount)
	if err != nil {
		t.Fatalf("Failed to check upload part presence: %v", err)
	}
	if partCount != 0 {
		t.Fatalf("Expected upload parts to be deleted, but they are still present")
	}
}

func TestDeleteUpload_NotFound(t *testing.T) {
	id := uuid.New()
	ctx, _, cleanup := SetupTest(t)
	defer cleanup()

	// Attempt to delete a non-existent upload
	err := DeleteUpload(ctx, id)
	var target ErrUploadNotFound
	if !errors.As(err, &target) {
		t.Fatalf("Expected ErrUploadNotFound, got %v", err)
	}
}

func TestGetUpload(t *testing.T) {
	id := uuid.New()
	ctx, _, cleanup := SetupTest(t)
	defer cleanup()

	// Create upload
	err := CreateUpload(ctx, id, 2, 1024, "application/octet-stream")
	if err != nil {
		t.Fatalf("Failed to create upload: %v", err)
	}

	// Get upload
	upload, err := GetUpload(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get upload: %v", err)
	}

	// Verify details
	if upload.ID != id {
		t.Fatalf("Expected upload ID %v, got %v", id, upload.ID)
	}
	if upload.PartsCount != 2 {
		t.Fatalf("Expected parts count 2, got %d", upload.PartsCount)
	}
	if upload.Size != 1024 {
		t.Fatalf("Expected size 1024, got %d", upload.Size)
	}
	if upload.MimeType != "application/octet-stream" {
		t.Fatalf("Expected mime type application/octet-stream, got %s", upload.MimeType)
	}
	if upload.Status != UploadStatusPending {
		t.Fatalf("Expected status %s, got %s", UploadStatusPending, upload.Status)
	}
	if upload.CreatedAt.IsZero() {
		t.Fatalf("Expected non-zero created_at, got %v", upload.CreatedAt)
	}
}

func TestGetUpload_NotFound(t *testing.T) {
	id := uuid.New()
	ctx, _, cleanup := SetupTest(t)
	defer cleanup()
	// Attempt to get a non-existent upload
	_, err := GetUpload(ctx, id)
	var target ErrUploadNotFound
	if !errors.As(err, &target) {
		t.Fatalf("Expected ErrUploadNotFound, got %v", err)
	}
}

func TestGetUploadParts(t *testing.T) {
	id := uuid.New()
	ctx, _, cleanup := SetupTest(t)
	defer cleanup()

	// Create upload
	err := CreateUpload(ctx, id, 6, 1024, "application/octet-stream")
	if err != nil {
		t.Fatalf("Failed to create upload: %v", err)
	}

	// Get upload parts
	parts, err := GetUploadParts(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get upload parts: %v", err)
	}

	// Verify parts details
	if len(parts) != 6 {
		t.Fatalf("Expected 6 parts, got %d", len(parts))
	}
}

func TestGetUploadParts_NotFound(t *testing.T) {
	id := uuid.New()
	ctx, _, cleanup := SetupTest(t)
	defer cleanup()

	// Attempt to get parts for a non-existent upload
	parts, err := GetUploadParts(ctx, id)
	var target ErrUploadNotFound
	if !errors.As(err, &target) {
		t.Fatalf("Expected ErrUploadNotFound, got Err: %v, parts: %v", err, parts)
	}
}
