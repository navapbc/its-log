package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	defaultstorage "github.com/jadudm/its-log/internal/default-storage"
	"github.com/spf13/viper"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob"
)

type BlobStorage struct {
	Counter int64
	Bucket  *blob.Bucket
}

// Memory buckets! great for testing.
// b, err := blob.OpenBucket(ctx, "mem://?prefix=a/subfolder/")

func (s *BlobStorage) Init() error {
	url := fmt.Sprintf("s3://%s?endpoint=%s://%s:%s&s3ForcePathStyle=true",
		viper.GetString("s3.bucket"),
		viper.GetString("s3.scheme"),
		viper.GetString("s3.host"),
		viper.GetString("s3.port"))

	fmt.Printf("OpenBucket: %s\n", url)
	b, err := blob.OpenBucket(context.Background(), url)
	if err != nil {
		log.Fatalf("Failed to setup bucket: %s", err)
	}

	s.Counter = 0
	s.Bucket = b
	return nil
}

func (s *BlobStorage) Event(source string, event string, value any, value_type string) (int64, error) {
	fmt.Printf("Blob %s %v %v\n", event, value, value_type)
	t := time.Now()
	secondsSinceEpoch := t.Unix()

	pfix := t.Format("2006-01-02")
	bucket := blob.PrefixedBucket(s.Bucket, pfix+"/")

	ctx := context.Background()
	object_name := fmt.Sprintf("%s-%d-%09d.json", pfix, secondsSinceEpoch, s.Counter)
	w, err := bucket.NewWriter(ctx, object_name, nil)
	if err != nil {
		return 0, err
	}

	evt := make(map[string]string)
	evt["source"] = source
	evt["event"] = event
	evt["value"] = fmt.Sprintf("%s", value)
	evt["type"] = value_type
	jsonString, err := json.MarshalIndent(evt, "", "  ")

	if err != nil {
		// TODO: handle error gracefully
		fmt.Println("error marshalling")
		d := &defaultstorage.DefaultStorage{}
		d.Event(source, event, value, value_type)
	}

	_, writeErr := fmt.Fprintln(w, string(jsonString))
	// Always check the return value of Close when writing.
	closeErr := w.Close()
	if writeErr != nil {
		log.Fatal(writeErr)
	}
	if closeErr != nil {
		log.Fatal(closeErr)
	}

	s.Counter += 1
	return s.Counter, nil
}

func (s *BlobStorage) Close() {
	s.Bucket.Close()
}
