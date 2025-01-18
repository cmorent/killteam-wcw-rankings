package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/cmorent/killteam-wcw-rankings/pkg/scoring"
)

type HandleInsertEventResultsRequest struct {
	Name     string   `json:"name"`
	Rankings []string `json:"rankings"`
}

func (s *Server) HandleInsertEventResults(w http.ResponseWriter, r *http.Request) {
	req := HandleInsertEventResultsRequest{}
	json.NewDecoder(r.Body).Decode(&req)

	err := s.db.InsertEventRankings(r.Context(), req.Name, req.Rankings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) HandleGetSeasonalRankings(w http.ResponseWriter, r *http.Request) {
	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		slog.Error(fmt.Sprintf("failed to parse year %q", r.URL.Query().Get("year")), slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if year == 0 {
		year, _, _ = time.Now().Date()
	}

	seasonalResults, err := s.db.GetSeasonEventsResults(r.Context(), year)
	if err != nil {
		slog.Error("failed to retrieve seasonal events rankings from db", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rankings, err := scoring.ComputeSeasonalRankings(seasonalResults)
	if err != nil {
		slog.Error("failed to compute seasonal rankings", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(rankings)
	if err != nil {
		slog.Error("failed to encode response", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
