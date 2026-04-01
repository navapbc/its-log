package etl

import (
	"context"
	"fmt"
	"log"

	"github.com/navapbc/its-log/internal/schema/models"
	"github.com/navapbc/its-log/internal/types"
)

func HashSummaries(etlP *types.RunEtlParams) error {

	summaryRows, err := etlP.Storage.Queries.GetAllSummaries(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching summaries: %s", err.Error())
	}

	for _, srow := range summaryRows {
		// This is the incoming hash from the the database
		currentHash := srow.Hash.String
		// recompute the hash
		srow.HashItslogSummary()
		if currentHash != srow.Hash.String {
			log.Printf("existingHash[%s] computedHash[%s]", currentHash, srow.Hash.String)
			err := etlP.Storage.Queries.InsertFullSummary(context.Background(), models.InsertFullSummaryParams{
				KeyID:     srow.KeyID,
				LastRun:   srow.LastRun,
				Date:      srow.Date,
				Operation: srow.Operation,
				Tags:      srow.Tags,
				Value:     srow.Value,
				Count:     srow.Count,
				Hash:      srow.Hash,
			})
			if err != nil {
				return fmt.Errorf("could not update hash for %s: %s", etlP.Storage.YYYYMMDD(), err.Error())
			}
		}
	}

	return nil

}
