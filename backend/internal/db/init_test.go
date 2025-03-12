package db

import (
	"context"
	"testing"
)

func TestInitDB(t *testing.T) {
	t.Parallel()
	// Call the InitDB function
	dbpool, cleanup, err := InitDBPool()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Ensure the cleanup function is called
	defer cleanup()

	// Check if the dbpool is not nil
	if dbpool == nil {
		t.Fatalf("Expected dbpool to be non-nil")
	}

	// Run a sample query to ensure the connection is working
	var result int
	err = dbpool.QueryRow(context.Background(), "SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("Expected no error from sample query, got %v", err)
	}

	// Check if the result is as expected
	if result != 1 {
		t.Fatalf("Expected result to be 1, got %v", result)
	}
}
