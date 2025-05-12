package main

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	cred, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		log.Fatalf("unable to retrieve AWS credentials, %v", err)
	}
	log.Printf("conf %#v %#v\n", cfg, cred)

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	r := strings.NewReader("some io.Reader stream to be read\n")

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("vknpmtbczsjxrqghlwadueoyfiwzkhqscilpxtnbvfjmr"),
		Key:    aws.String("someiofromgo"),
		Body:   r,
	})

	if err != nil {
		log.Fatalf("unable to upload object, %v", err)
	}
	log.Printf("%#v\n", result)
}
