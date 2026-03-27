package e2e

import (
	"log"
	"strings"

	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/types"
)

func RunDefaultSequence() *types.Storage {
	apiKey, err := base.GetKeyBundle("pupper", "pup_admin")
	if err != nil {
		panic("could not find key for pupper/pup_admin")
	}
	s := types.NewStorage(apiKey.AppId)
	s.Init()
	today := s.YYYYMMDD()
	target := "/v1" + constants.SEQUENCE_RUN
	target = strings.Replace(target, ":date", today, -1)
	target = strings.Replace(target, ":name", "default", -1)
	log.Println("running default sequence")
	get(buildBase()+target, apiKey)
	return s
}
