package storage

import (
	"context"
	"strconv"
	"sync"
	"testing"
)

func TestUploadAndDownload(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Upload object
	objectName := t.Name() + "-test-object"
	uploadedData := []byte("This is a test object")
	err := Upload(ctx, objectName, uploadedData)
	if err != nil {
		t.Fatalf("Failed to upload object: %v", err)
	}

	// Download object
	downloadedData, err := Download(ctx, objectName)
	if err != nil {
		t.Fatalf("Failed to download object: %v", err)
	}

	// Assert downloaded object = uploaded object
	if string(downloadedData) != string(uploadedData) {
		t.Fatalf("Downloaded data does not match uploaded data. Got: %s, Want: %s", string(downloadedData), string(uploadedData))
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Upload object
	objectName := t.Name() + "-test-object"
	uploadedData := []byte("This object will be deleted")
	err := Upload(ctx, objectName, uploadedData)
	if err != nil {
		t.Fatalf("Failed to upload object: %v", err)
	}

	// Delete object
	err = Delete(ctx, objectName)
	if err != nil {
		t.Fatalf("Failed to delete object: %v", err)
	}

	// Try to download deleted object
	_, err = Download(ctx, objectName)
	if err == nil {
		t.Fatalf("Expected error when downloading deleted object, but got none")
	}

}

func TestExists(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	objectName := t.Name() + "-test-object"

	// Try deleting the object
	Delete(ctx, objectName)

	// Initially, the object should not exist
	exists, err := Exists(ctx, objectName)
	if err != nil {
		t.Fatalf("Failed to check if object exists: %v", err)
	}
	if exists {
		t.Fatalf("Object should not exist initially")
	}

	// Upload object
	uploadedData := []byte("This is a test object")
	err = Upload(ctx, objectName, uploadedData)
	if err != nil {
		t.Fatalf("Failed to upload object: %v", err)
	}

	// Now the object should exist
	exists, err = Exists(ctx, objectName)
	if err != nil {
		t.Fatalf("Failed to check if object exists: %v", err)
	}
	if !exists {
		t.Fatalf("Object should exist after upload")
	}
}

func TestCompose(t *testing.T) {
	ctx := context.Background()

	t.Run("0 source object", func(t *testing.T) {
		t.Parallel()
		dstObjectName := t.Name() + "-dst-object"
		err := Compose(ctx, dstObjectName, []string{})
		if err == nil {
			t.Fatalf("Expected error when composing with 0 source objects, but got none")
		}
	})

	t.Run("1 source object", func(t *testing.T) {
		t.Parallel()
		dstObjectName := t.Name() + "-dst-object"
		srcObjectName := t.Name() + "-src-object-1"
		err := Upload(ctx, srcObjectName, []byte("This is a test object"))
		if err != nil {
			t.Fatalf("Failed to upload source object: %v", err)
		}
		err = Compose(ctx, dstObjectName, []string{srcObjectName})
		if err == nil {
			t.Fatalf("Expected error when composing with 1 source object, but got none")
		}
	})

	t.Run("5 source objects", func(t *testing.T) {
		t.Parallel()
		dstObjectName := t.Name() + "-dst-object"
		srcObjectNames := []string{}
		expectedData := []byte{}

		var wg sync.WaitGroup
		var mu sync.Mutex
		for i := 1; i <= 5; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				srcObjectName := t.Name() + "-src-object-" + strconv.Itoa(i)
				data := []byte("This is test object " + strconv.Itoa(i))
				err := Upload(ctx, srcObjectName, data)
				if err != nil {
					t.Errorf("Failed to upload source object %d: %v", i, err)
				}
				mu.Lock()
				srcObjectNames = append(srcObjectNames, srcObjectName)
				expectedData = append(expectedData, data...)
				mu.Unlock()
			}(i)
		}
		wg.Wait()
		err := Compose(ctx, dstObjectName, srcObjectNames)
		if err != nil {
			t.Fatalf("Failed to compose objects: %v", err)
		}
		composedData, err := Download(ctx, dstObjectName)
		if err != nil {
			t.Fatalf("Failed to download composed object: %v", err)
		}
		if string(composedData) != string(expectedData) {
			t.Fatalf("Composed object data does not match expected data")
		}
	})

	t.Run("33 source objects", func(t *testing.T) {
		t.Parallel()
		dstObjectName := t.Name() + "-dst-object"
		srcObjectNames := []string{}
		var wg sync.WaitGroup
		var mu sync.Mutex
		for i := 1; i <= 33; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				srcObjectName := t.Name() + "-src-object-" + strconv.Itoa(i)
				data := []byte("This is test object " + strconv.Itoa(i))
				err := Upload(ctx, srcObjectName, data)
				if err != nil {
					t.Errorf("Failed to upload source object %d: %v", i, err)
				}
				mu.Lock()
				srcObjectNames = append(srcObjectNames, srcObjectName)
				mu.Unlock()
			}(i)
		}
		wg.Wait()
		err := Compose(ctx, dstObjectName, srcObjectNames)
		if err == nil {
			t.Fatalf("Expected error when composing with 33 source objects, but got none")
		}
	})

}
