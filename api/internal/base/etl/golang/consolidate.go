package etl

import (
	"context"
	"fmt"
	"log"

	"github.com/navapbc/its-log/internal/base"
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
		"prior_summary_dates_to_include",
	}

	err := hasExpectedKeys("consolidate", etlP, expectedConsolidateKeys)
	if err != nil {
		return err
	}

	// dateToConsolidateTo := etlP.Storage.YYYYMMDD()

	prior_dates_any := etlP.Payload["prior_summary_dates_to_include"]
	// For each of those keys, let's check the date, and if it parses, look to see if there's a DB we can process.
	for _, past_date_any := range prior_dates_any.([]any) {
		if past_date, ok := past_date_any.(string); ok {
			appId := base.GetOrPanic(etlP.GinCtx, "AppId")
			keyId := base.GetOrPanic(etlP.GinCtx, "KeyId")

			// Init the past storage
			past_storage := types.NewStorage(appId)
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
			// current storage, which is in the etlP parameter.
			past_storage.Lock()
			defer past_storage.Unlock()

			summaryRows, err := past_storage.Queries.GetAllSummaries(context.Background())
			if err != nil {
				return fmt.Errorf("error fetching prior summaries: %s", err.Error())
			}

			for _, srow := range summaryRows {
				// type InsertSummaryParams struct {
				// 	Date      string
				// 	Operation string
				// 	Tags      string
				// 	Value     string
				// }
				etlP.Storage.Lock()

				err := etlP.Storage.Queries.InsertFullSummary(context.Background(), models.InsertFullSummaryParams{
					KeyID:     keyId,
					LastRun:   srow.LastRun,
					Date:      srow.Date,
					Operation: srow.Operation,
					Tags:      srow.Tags,
					Value:     srow.Value,
					Count:     srow.Count,
				})
				if err != nil {
					etlP.Storage.Unlock()
					return fmt.Errorf("could not insert summary from %s: %s", past_date, err.Error())
				}
				etlP.Storage.Unlock()
			}

			// Close things so writes complete.
			// Deferred locks will be released.
			past_storage.Close()
		} else {
			return fmt.Errorf("could not convert: %v", prior_dates_any)
		}
	}

	return nil
}
