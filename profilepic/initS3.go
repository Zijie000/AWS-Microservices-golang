package profilepic

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client
var bucketName = ""

func init() {
	bucketName = os.Getenv("BUCKET_NAME")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		panic("Cannot load aws config: " + err.Error())
	}
	s3Client = s3.NewFromConfig(cfg)
}
