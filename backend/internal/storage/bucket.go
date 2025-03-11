package storage

import (
	"context"
	"errors"
	"io"

	"cloud.google.com/go/storage"
)

func Download(ctx context.Context, objectName string) ([]byte, error) {
	bucket, err := Bucket()
	if err != nil {
		return nil, err
	}
	object := bucket.Object(objectName)
	reader, err := object.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Upload(ctx context.Context, objectName string, data []byte) error {
	bucket, err := Bucket()
	if err != nil {
		return err
	}
	object := bucket.Object(objectName)
	writer := object.NewWriter(ctx)
	writer.Write(data)
	return writer.Close()
}

func Delete(ctx context.Context, objectName string) error {
	bucket, err := Bucket()
	if err != nil {
		return err
	}
	object := bucket.Object(objectName)
	return object.Delete(ctx)
}

func Exists(ctx context.Context, objectName string) (bool, error) {
	bucket, err := Bucket()
	if err != nil {
		return false, err
	}
	object := bucket.Object(objectName)
	_, err = object.NewReader(ctx)
	if errors.Is(err, storage.ErrObjectNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func Compose(ctx context.Context, dstObjectName string, srcObjectNames []string) error {
	if len(srcObjectNames) == 0 {
		return errors.New("no object names provided")
	}
	if len(srcObjectNames) == 1 {
		return errors.New("only one object name provided")
	}
	if len(srcObjectNames) > 32 {
		return errors.New("only up to 32 object names can be provided")
	}
	bucket, err := Bucket()
	if err != nil {
		return err
	}
	dstObject := bucket.Object(dstObjectName)
	srcObjects := make([]*storage.ObjectHandle, len(srcObjectNames))
	for i, srcObjectName := range srcObjectNames {
		srcObjects[i] = bucket.Object(srcObjectName)
	}
	composer := dstObject.ComposerFrom(srcObjects...)
	_, err = composer.Run(ctx)
	return err
}
