package e2e

import (
	"runtime"
	"strings"
	"testing"

	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/etl"
	"github.com/navapbc/its-log/internal/types"
)

func RunSequence(t *testing.T, dateOffset int, sequenceName string) *types.Storage {
	date := types.NewILTimeToday()
	date.SubtractDays(dateOffset)

	apiKey, err := base.GetKeyBundle("pupper", "pup_admin")
	if err != nil {
		panic("could not find key for pupper/pup_admin")
	}
	s, err := types.NewStorage(apiKey.AppId)
	if err != nil {
		t.Error(err)
	}
	s.SetDateILT(date)
	err = s.InitDB()
	if err != nil {
		t.Error(err)
	}

	// DEBUG LOG
	// log.Printf("RunSequence: %s <- %s\n", s.ILTime.AsYYYYMMDD(), s.Filename)

	// If we don't do an insert, then nothing loads the
	// default ETLs. For testing, we have to force the issue.
	pc, _, _, _ := runtime.Caller(0)
	funcName := runtime.FuncForPC(pc).Name()
	etl.LoadDefaultEtlFiles(s, funcName)

	target := "/v1" + constants.SEQUENCE_RUN
	target = strings.Replace(target, ":date", date.AsYYYYMMDD(), -1)
	target = strings.Replace(target, ":name", sequenceName, -1)

	// DEBUG LOG
	// log.Println("running sequence: " + sequenceName)
	// log.Printf("sequence target date: %s\n", date.AsYYYYMMDD())

	bundle := make(map[string]any)
	bundle["dates-to-backup"] = []string{date.AsYYYYMMDD()}
	res := post(buildBase()+target, bundle, apiKey)
	if !checkRes(res) {
		t.Fail()
	}
	return s
}
