package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
)

type Storage struct {
	client     *storage.Client
	bucketName string
}

func New(ctx context.Context, bucketName string) (*Storage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to init GCP storage client: %w", err)
	}

	return &Storage{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (db *Storage) InsertEventRankings(ctx context.Context, eventName string, rankings []string) error {
	year, _, _ := time.Now().Date()
	dbName := fmt.Sprintf("%d.json", year)

	eventsRankings, err := db.downloadDB(ctx, dbName)
	if err != nil {
		return fmt.Errorf("failed to download db file: %w", err)
	}

	eventsRankings[eventName] = rankings

	err = db.uploadDB(ctx, dbName, eventsRankings)
	if err != nil {
		return fmt.Errorf("failed to upload db file: %w", err)
	}
	return nil
}

func (db *Storage) GetSeasonEventsResults(ctx context.Context, year int) (map[string][]string, error) {
	dbName := fmt.Sprintf("%d.json", year)

	eventsRankings, err := db.downloadDB(ctx, dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to download db file: %w", err)
	}

	return eventsRankings, nil
}

func (db *Storage) downloadDB(ctx context.Context, dbName string) (map[string][]string, error) {
	sr, err := db.client.Bucket(db.bucketName).Object(dbName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader for %q: %w", fmt.Sprintf("%s/%s", db.bucketName, dbName), err)
	}
	defer sr.Close()

	var jsonDB map[string][]string
	err = json.NewDecoder(sr).Decode(&jsonDB)
	if err != nil {
		return nil, fmt.Errorf("failed to read GCS object: %w", err)
	}

	return jsonDB, nil
}

func (db *Storage) uploadDB(ctx context.Context, dbName string, jsonDB map[string][]string) error {
	sw := db.client.Bucket(db.bucketName).Object(dbName).NewWriter(ctx)
	defer sw.Close()

	err := json.NewEncoder(sw).Encode(jsonDB)
	if err != nil {
		return fmt.Errorf("failed to upload db: %w", err)
	}
	return nil
}

func (db *Storage) Shutdown(ctx context.Context) error {
	return db.client.Close()
}
