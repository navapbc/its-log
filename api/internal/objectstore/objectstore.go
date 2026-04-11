package objectstore

import (
	"context"
	"fmt"
	"log"

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

// FIXME: This can be broken into several setters, and then an .Init() method
func NewObjectStore(provider string) *ObjectStore {
	os := &ObjectStore{
		Provider: provider,
	}
	return os
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

func (os *ObjectStore) Init() {
	client, err := NewObjectClient(os)
	if err != nil {
		panic("could not initialize object store")
	}
	os.ObjectClient = client
}

func (os *ObjectStore) ConstructUrl() string {
	return fmt.Sprintf("s3://%s", os.Bucket)
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

func (os *ObjectStore) Write(obj []byte) {
	ctx := context.Background()
	bucket, err := s3blob.OpenBucketV2(ctx, os.ObjectClient.AwsS3Client, os.Bucket, nil)

	if err != nil {
		log.Printf("could not open bucket: %v", err)
	}
	defer bucket.Close()

	// Open the key "foo.txt" for writing with the default options.
	w, err := bucket.NewWriter(ctx, "hello.txt", nil)
	if err != nil {
		log.Println("bucket error: " + err.Error())
	}
	_, writeErr := fmt.Fprintln(w, obj)
	// Always check the return value of Close when writing.
	closeErr := w.Close()
	if writeErr != nil {
		log.Println(writeErr)
	}
	if closeErr != nil {
		log.Println(closeErr)
	}
}
