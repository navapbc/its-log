package etl

import (
	"context"
	"fmt"
	"log"

	"github.com/navapbc/its-log/internal/schema/models"
	"github.com/navapbc/its-log/internal/types"
)

func Consolidate(etlP *types.RunEtlParams) error {

	// Consolidation assumes we want to pull prior summary table(s) into the
	// current summary table. We rely on the uniqueness constraints and INSERT OR REPLACE
	// as a way to make sure we don't get duplicates.
	//
	// If we do this daily, we end up with the current table always having everything from the summary history.
	// Or, today's summary includes all prior summaries.
	// If we miss a day, we can always grab multiple prior DBs, and (again) rely on uniqueness.
	//
	// This will use the extended keys feature to allow us to specify the date(s) we want to pull
	// forward into the current table.
	var expectedConsolidateKeys = []string{
		"prior-summary-dates-to-include",
	}

	err := hasExpectedKeys("consolidate", etlP, expectedConsolidateKeys)
	if err != nil {
		return err
	}

	prior_dates_any := etlP.Payload["prior_summary_dates_to_include"]
	if prior_dates_any == nil {
		prior_dates_any = make([]any, 0)
	}
	for _, past_date_any := range prior_dates_any.([]any) {
		// For each of those keys, let's check the date, and if it parses, look to see if there's a DB we can process.
		if past_date, ok := past_date_any.(string); ok {

			// Init the past storage
			past_storage := types.NewStorage(etlP.AppId)
			err := past_storage.SetDate(past_date)
			if err != nil {
				return fmt.Errorf("could not parse date; must be YYYY-MM-DD: %s", past_date)
			}
			err = past_storage.Init()
			if err != nil {
				log.Println("past storage init error: " + err.Error())
				panic(err)
			}

			// Now, we have past storage set up. Copy all of the summary table from it to the
			// current storage, which is in the etlP parameter. Locking here is possibly
			// not *strictly* necessary, but safer than not?
			past_storage.Lock()
			defer past_storage.Unlock()

			summaryRows, err := past_storage.Queries.GetAllSummaries(context.Background())
			if err != nil {
				return fmt.Errorf("error fetching prior summaries: %s", err.Error())
			}

			// etlP.Storage.Lock()

			for _, srow := range summaryRows {

				err := etlP.Storage.Queries.InsertFullSummary(context.Background(), models.InsertFullSummaryParams{
					KeyID:     etlP.KeyId,
					LastRun:   srow.LastRun,
					Date:      srow.Date,
					Operation: srow.Operation,
					Tags:      srow.Tags,
					Value:     srow.Value,
					Count:     srow.Count,
				})
				if err != nil {
					// etlP.Storage.Unlock()
					return fmt.Errorf("could not insert summary from %s: %s", past_date, err.Error())
				}
			}

			// Deferred locks will be released.
			past_storage.Close()
			// etlP.Storage.Unlock()
		} else {
			return fmt.Errorf("could not convert: %v", prior_dates_any)
		}
	}

	return nil
}
