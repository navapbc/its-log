package defaultstorage

import (
	"fmt"
	"log"
	"time"

	"github.com/jadudm/its-log/internal/itslog"
	_ "gocloud.dev/blob/s3blob"
)

type DefaultStorage struct {
	Counter int64
}

// Memory buckets! great for testing.
// b, err := blob.OpenBucket(ctx, "mem://?prefix=a/subfolder/")

func (s *DefaultStorage) Init() error {
	fmt.Println("default storage initialized")
	s.Counter = 0
	return nil
}

func (s *DefaultStorage) Event(e *itslog.Event) (int64, error) {
	t := time.Now()
	pfix := t.Format("2006-01-02")
	log.Printf("%06d %s %s %s\n", s.Counter, pfix, e.Source, e.Event)

	s.Counter += 1
	return s.Counter, nil
}
