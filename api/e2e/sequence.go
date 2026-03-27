package e2e

import (
	"fmt"

	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/types"
)

func RunDefaultSequence() {
	apiKey, err := base.GetKeyBundle("pupper", "pup_admin")
	if err != nil {
		panic("could not find key for pupper/pup_admin")
	}
	s := types.NewStorage(apiKey.AppId)
	today := s.YYYYMMDD()
	targetUrl := fmt.Sprintf("%s/v1/sequence/%s/default", buildBase(), today)
	get(targetUrl, apiKey)
}
