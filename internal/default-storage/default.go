package defaultstorage

import (
	"fmt"
	"time"

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

func (s *DefaultStorage) Event(source string, event string, value any, value_type string) (int64, error) {
	t := time.Now()
	pfix := t.Format("2006-01-02")
	fmt.Printf("%06d %s %s %s %s %s\n", s.Counter, pfix, source, event, value, value_type)

	s.Counter += 1
	return s.Counter, nil
}
