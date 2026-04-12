package etl

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/navapbc/its-log/internal/objectstore"
	"github.com/navapbc/its-log/internal/types"
	"github.com/spf13/viper"
)

func copyToBackupBucket(s *types.Storage) (int, error) {
	os := objectstore.NewObjectStore("aws")
	os.SetBucket(viper.GetString("backup.bucket"))
	os.SetEndpoint(
		viper.GetString("backup.aws.endpoint.scheme"),
		viper.GetString("backup.aws.endpoint.host"),
		viper.GetString("backup.aws.endpoint.port"))
	os.SetRegion(viper.GetString("backup.aws.region"))
	os.SetKeyIdAndAccessKey(
		viper.GetString("backup.aws.keyid"),
		viper.GetString("backup.aws.accesskey"))
	os.Init()
	// os.WriteBytes([]string{"here", "there"}, []byte("hello"))

	dest := []string{time.Now().UTC().Format("2006-01-02"), s.Filename}

	log.Println("source: " + strings.Join(s.Path, "/"))
	log.Println("dest: " + strings.Join(dest, "/"))

	written, err := os.CopyToS3(s.Path, dest)
	if err != nil {
		return -1, err
	}
	return written, nil
}

func Backup(etlP *types.RunEtlParams) error {
	var expectedConsolidateKeys = []string{
		"dates-to-backup",
	}

	err := hasExpectedKeys("backup", etlP, expectedConsolidateKeys)
	if err != nil {
		return err
	}

	for _, d := range etlP.Payload["dates-to-backup"].([]any) {
		// We want to copy a database file at a given location.
		// Create a storage object, and close it.
		date := d.(string)
		storage := types.NewStorage(etlP.AppId)
		err := storage.SetDateYMD(date)
		if err != nil {
			return fmt.Errorf("could not parse date; must be YYYY-MM-DD: %s", date)
		}
		err = storage.Init()
		if err != nil {
			log.Println("storage init error: " + err.Error())
			panic(err)
		}

		// Lock the DB while copying.
		storage.Lock()
		storage.Close()
		written, err := copyToBackupBucket(storage)
		storage.Unlock()
		if err != nil {
			return err
		}
		log.Println("wrote bytes: " + strconv.Itoa(written))
	}
	return nil
}
