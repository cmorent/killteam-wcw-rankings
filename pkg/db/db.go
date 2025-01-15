package db

import "context"

type DB interface {
	InsertEventRankings(ctx context.Context, name string, rankings []string) error
	GetSeasonEventsResults(ctx context.Context, year int) (map[string][]string, error)
	Shutdown(ctx context.Context) error
}
