package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/cmorent/killteam-wcw-rankings/pkg/db"
	"github.com/rs/cors"
)

type Server struct {
	*http.Server
	db db.DB
}

type HandleInsertEventResultsRequest struct {
	Name     string   `json:"name"`
	Rankings []string `json:"rankings"`
}

func New(addr string, db db.DB) (*Server, error) {
	s := &Server{
		Server: &http.Server{
			Addr: addr,
		},
		db: db,
	}

	m := http.NewServeMux()
	m.HandleFunc("POST /events", s.HandleInsertEventResults)
	m.HandleFunc("GET /rankings", s.HandleGetSeasonalRankings)
	s.Handler = cors.AllowAll().Handler(m)
	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	done := make(chan error)

	go func() {
		slog.Info("server now listening on " + s.Addr)
		done <- s.ListenAndServe()
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return s.Shutdown(context.Background())
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.db.Shutdown(ctx); err != nil {
		slog.Error("failed to close gcp storage client")
	}
	return nil
}
