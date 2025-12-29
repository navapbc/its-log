package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jadudm/its-log/internal/sqlite/models"
)

func sumBySource(s *SqliteStorage) {
	ctx := context.Background()
	rows, err := s.queries.SourceCountsForTheDay(ctx)
	if err != nil {
		// TODO error handling
		panic(err)
	}

	for _, r := range rows {
		source_name, _ := s.queries.GetSourceName(ctx, r.Source)
		err = s.queries.InsertSummary(ctx, models.InsertSummaryParams{
			Operation: "count.day.by_source",
			Source:    source_name,
			Event:     sql.NullString{String: "", Valid: false},
			Value:     float64(r.SourceCount),
		})
		if err != nil {
			panic(err)
		}
	}
}

func sumByEvent(s *SqliteStorage) {
	ctx := context.Background()
	rows, err := s.queries.EventCountsForTheDay(ctx)
	if err != nil {
		// TODO error handling
		panic(err)
	}

	for _, r := range rows {
		source_name, _ := s.queries.GetSourceName(ctx, r.Source)
		event_name, _ := s.queries.GetEventName(ctx, models.GetEventNameParams{
			SourceHash: r.Source,
			EventHash:  r.Event,
		})

		err = s.queries.InsertSummary(ctx, models.InsertSummaryParams{
			Operation: "count.day.by_event",
			Source:    source_name,
			Event:     sql.NullString{String: event_name, Valid: true},
			Value:     float64(r.EventCount),
		})
		if err != nil {
			panic(err)
		}
	}
}

func eventsByHour(s *SqliteStorage) {
	ctx := context.Background()
	rows, err := s.queries.EventCountsByTheHour(ctx)
	if err != nil {
		// TODO error handling
		panic(err)
	}

	for _, r := range rows {
		source_name, _ := s.queries.GetSourceName(ctx, r.Source)
		event_name, _ := s.queries.GetEventName(ctx, models.GetEventNameParams{
			SourceHash: r.Source,
			EventHash:  r.Event,
		})

		err = s.queries.InsertSummary(ctx, models.InsertSummaryParams{
			Operation: fmt.Sprintf("count.source_and_event.by_hour.%s", r.Hour),
			Source:    source_name,
			Event:     sql.NullString{String: event_name, Valid: true},
			Value:     float64(r.EventCount),
		})
		if err != nil {
			panic(err)
		}
	}
}

func sourceByHour(s *SqliteStorage) {
	ctx := context.Background()
	rows, err := s.queries.SourceCountsByTheHour(ctx)
	if err != nil {
		// TODO error handling
		panic(err)
	}

	for _, r := range rows {
		source_name, _ := s.queries.GetSourceName(ctx, r.Source)

		err = s.queries.InsertSummary(ctx, models.InsertSummaryParams{
			Operation: fmt.Sprintf("count.source.by_hour.%s", r.Hour),
			Source:    source_name,
			Event:     sql.NullString{String: "", Valid: false},
			Value:     float64(r.SourceCount),
		})
		if err != nil {
			panic(err)
		}
	}
}

func (s *SqliteStorage) Summarize() {
	sumBySource(s)
	sumByEvent(s)
	eventsByHour(s)
	sourceByHour(s)
}
