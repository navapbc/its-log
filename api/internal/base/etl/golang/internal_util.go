package etl

import (
	"fmt"

	"github.com/navapbc/its-log/internal/types"
)

func hasExpectedKeys(stepName string, etlP *types.RunEtlParams, jsonKeyList []string) error {
	for _, key := range jsonKeyList {
		_, ok := etlP.Payload[key]
		if !ok {
			return fmt.Errorf("missing parameter %s for %s", key, stepName)
		}
	}
	return nil
}
