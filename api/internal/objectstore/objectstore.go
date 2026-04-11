package objectstore

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"gocloud.dev/blob/s3blob"
	_ "gocloud.dev/blob/s3blob"
)

type ObjectKeys struct {
	KeyID  string
	Secret string
}

type ObjectClient struct {
	AwsS3Config *aws.Config
	AwsS3Client *s3.Client

	// Other providers can be added as needed.
}

type ObjectStore struct {
	Provider string

	ObjectClient   *ObjectClient
	EndpointScheme string
	EndpointHost   string
	EndpointPort   string
	Region         string
	Bucket         string
	Path           string
	Keys           *ObjectKeys
}

func NewObjectClient(os *ObjectStore) (*ObjectClient, error) {
	oclient := &ObjectClient{}
	switch os.Provider {
	case "aws":
		creds := credentials.NewStaticCredentialsProvider(os.Keys.KeyID, os.Keys.Secret, "")

		cfg, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithCredentialsProvider(creds),
			config.WithRegion(os.Region))
		if err != nil {
			log.Println("error with default config: " + err.Error())
			return nil, err
		}
		oclient.AwsS3Config = &cfg

		if os.EndpointHost != "" {
			log.Println("specifying custom endpoint: " + os.EndpointHost)

			oclient.AwsS3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
				o.BaseEndpoint = aws.String(fmt.Sprintf(
					"%s://%s:%s",
					os.EndpointScheme,
					os.EndpointHost,
					os.EndpointPort,
				))
				// FIXME: May work for ministack, but possibly not in production
				o.UsePathStyle = true
			})
		} else {
			log.Println("specifying default endpoint: " + os.EndpointHost)
			oclient.AwsS3Client = s3.NewFromConfig(cfg)
		}

	}
	return oclient, nil
}

func NewObjectStore(provider string) *ObjectStore {
	os := &ObjectStore{
		Provider: provider,
	}
	return os
}

func (os *ObjectStore) Init() {
	client, err := NewObjectClient(os)
	if err != nil {
		panic("could not initialize object store")
	}
	os.ObjectClient = client
}

func (os *ObjectStore) SetRegion(r string) *ObjectStore {
	os.Region = r
	return os
}

func (os *ObjectStore) SetKeyIdAndAccessKey(k string, s string) *ObjectStore {
	os.Keys = &ObjectKeys{
		KeyID:  k,
		Secret: s,
	}
	return os
}

func (os *ObjectStore) SetBucket(b string) *ObjectStore {
	os.Bucket = b
	return os
}

func (os *ObjectStore) SetPath(path string) *ObjectStore {
	os.Path = path
	return os
}

func (os *ObjectStore) SetEndpoint(scheme string, host string, port string) *ObjectStore {
	os.EndpointScheme = scheme
	os.EndpointHost = host
	os.EndpointPort = port
	return os
}

func (os *ObjectStore) ConstructUrl() string {
	return fmt.Sprintf("s3://%s", os.Bucket)
}

func (os *ObjectStore) WriteBytes(s3Path []string, obj []byte) (int, error) {
	r := bytes.NewReader(obj)
	written, err := os.Write(s3Path, r)
	return written, err
}

func (os *ObjectStore) Write(s3Path []string, reader io.Reader) (int, error) {
	ctx := context.Background()

	bucket, bucketErr := s3blob.OpenBucketV2(ctx, os.ObjectClient.AwsS3Client, os.Bucket, nil)
	if bucketErr != nil {
		log.Printf("could not open bucket: %v", bucketErr)
		return -1, bucketErr
	}
	defer bucket.Close()

	fullPath := strings.Join(s3Path, "/")
	// NOTE: The options to NewWriter could be of use (e.g. conditional writes)

	w, newWriterErr := bucket.NewWriter(ctx, fullPath, nil)
	if newWriterErr != nil {
		log.Println("writer creation error: " + newWriterErr.Error())
		return -1, newWriterErr
	}

	srcBuff := bufio.NewReader(reader)
	bytesWritten, writeErr := io.Copy(w, srcBuff)

	if writeErr != nil {
		log.Println("write err: " + writeErr.Error())
		if w != nil {
			w.Close()
		}
		return -1, writeErr
	}
	// Always check the return value of Close when writing.
	closeErr := w.Close()
	if closeErr != nil {
		log.Println("writer close err: " + closeErr.Error())
		return -1, closeErr
	}

	return int(bytesWritten), nil
}

func (objs *ObjectStore) CopyToS3(sourcePath []string, destPath []string) (int, error) {
	source := filepath.Join(sourcePath...)
	if _, err := os.Stat(source); errors.Is(err, os.ErrNotExist) {
		log.Println("path not found locally: " + err.Error())
		return -1, err
	}

	srcReader, openErr := os.Open(source)
	if openErr != nil {
		log.Println("file open err: " + openErr.Error())
		return -1, openErr
	}
	srcBuff := bufio.NewReader(srcReader)

	bytesWritten, writeErr := objs.Write(destPath, srcBuff)
	if writeErr != nil {
		return -1, writeErr
	}

	return bytesWritten, nil
}
