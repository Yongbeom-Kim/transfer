package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateUpload(t *testing.T) {
	// Initialize the database connection
	dbpool, cleanup, err := InitDBPool()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer cleanup()
	tx, err := dbpool.Begin(context.Background())
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())
	ctx := WithTx(context.Background(), &tx)

	id := uuid.New()

	err = CreateUpload(ctx, id, 1, 1024, "image/jpeg")
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
	err = tx.QueryRow(context.Background(), "(SELECT * FROM upload.uploads WHERE id = $1)", id).Scan(&upload.ID, &upload.CreatedAt, &upload.Status, &upload.PartsCount, &upload.Size, &upload.MimeType)
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
	err = tx.QueryRow(context.Background(), "SELECT COUNT(*) FROM upload.parts WHERE upload_id = $1", id).Scan(&partsCount)
	if err != nil {
		t.Fatalf("Failed to check if upload parts exist: %v", err)
	}
	if partsCount != 1 {
		t.Fatalf("Expected 1 upload part, got %d", partsCount)
	}

}
