package e2e

import (
	"log"
	"strings"
	"time"

	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/types"
)

func RunDefaultSequence(dateOffset int) *types.Storage {
	date := time.Now().AddDate(0, 0, dateOffset).Format("2006-01-02")
	log.Println("running for date: " + date)

	apiKey, err := base.GetKeyBundle("pupper", "pup_admin")
	if err != nil {
		panic("could not find key for pupper/pup_admin")
	}
	s := types.NewStorage(apiKey.AppId)
	s.SetDate(date)
	s.Init()
	// If we don't do an insert, then nothing loads the
	// default ETLs. For testing, we have to force the issue.
	base.LoadDefaultEtlFiles(s)

	target := "/v1" + constants.SEQUENCE_RUN
	target = strings.Replace(target, ":date", date, -1)
	target = strings.Replace(target, ":name", "default", -1)
	log.Println("running default sequence")
	post(buildBase()+target, make(map[string]any), apiKey)
	return s
}
