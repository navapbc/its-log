package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
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

// This wants a better design.
// And, it will want to pull from the summaries, not the raw data.
func renderBarCharts(s *SqliteStorage) {
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "events by hour",
		Subtitle: "separated by source",
	}))

	hours := make([]string, 24)
	for n := range 24 {
		hours = append(hours, fmt.Sprintf("%02d", n))
	}
	ctx := context.Background()
	names, err := s.queries.GetSourceNames(ctx)
	if err != nil {
		panic(err)
	}

	bar.SetXAxis(names)
	val_arr := make([]opts.BarData, len(names))
	events_for_source, _ := s.queries.SourceCountsForTheDay(ctx)
	for _, row := range events_for_source {
		val_arr = append(val_arr, opts.BarData{
			Value: row.SourceCount,
		})
	}
	bar.AddSeries("Source", val_arr)

	f, _ := os.Create("bar.html")
	bar.Render(f)
}

func (s *SqliteStorage) Summarize() {
	sumBySource(s)
	sumByEvent(s)
	eventsByHour(s)
	sourceByHour(s)
	renderBarCharts(s)
}
