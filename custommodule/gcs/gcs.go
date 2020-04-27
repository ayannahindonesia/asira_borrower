package gcs

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

// WriteFile writes file to gcs
func WriteFile(fPath string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	f, err := os.Open(fPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Sets your Google Cloud Platform project ID.
	// projectID := os.Getenv("GOOGLE_PROJECT_ID")

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	wc := client.Bucket(os.Getenv("GCS_BUCKET")).Object(fPath).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}
