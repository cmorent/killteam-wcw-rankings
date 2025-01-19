package db

import "context"

type DB interface {
	InsertEventRankings(ctx context.Context, name string, rankings map[string]int) error
	GetSeasonEventsResults(ctx context.Context, year int) (map[string]map[string]int, error)
	Shutdown(ctx context.Context) error
}
