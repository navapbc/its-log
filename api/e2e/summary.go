package e2e

import (
	"log"

	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/types"
)

func CheckSummaryValue(operation string, tags string, expected int, date *types.ILTime) bool {

	apiKey, err := base.GetKeyBundle("pupper", "pup_admin")
	if err != nil {
		panic("could not find key for pupper/pup_admin")
	}
	s := types.NewStorage(apiKey.AppId)
	s.SetDateILT(date)
	s.Init()
	target := "/v1" + constants.SUMMARY_READ

	bundle := map[string]any{
		"tags":      tags,
		"operation": operation,
		"date":      date.AsYYYYMMDD(),
	}
	mapResponse := post(buildBase()+target, bundle, apiKey)

	v, ok := mapResponse["count"]
	if ok {
		count := int(v.(float64))
		if expected == count {
			log.Printf("✅  %s: %d\n", operation, count)
			return true
		} else {
			log.Printf("⚠️  %s: expected %d, found %d\n", operation, expected, mapResponse["count"])
			return false
		}
	} else {
		log.Printf("🛑  %s: nil\n", operation)
		return false
	}
}
