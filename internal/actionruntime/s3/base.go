package s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/mitchellh/mapstructure"
)

func (s *Connector) getConnectionWithOptions(resourceOptions map[string]interface{}) (*s3.Client, error) {
	if err := mapstructure.Decode(resourceOptions, &s.ResourceOpts); err != nil {
		return nil, err
	}

	// format the parameters for the session you want to create.
	creds := credentials.NewStaticCredentialsProvider(s.ResourceOpts.AccessKeyID, s.ResourceOpts.SecretAccessKey, "")

	var cfg aws.Config
	var err error
	if s.ResourceOpts.Endpoint {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: s.ResourceOpts.BaseURL,
			}, nil
		})
		cfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(s.ResourceOpts.Region),
			config.WithCredentialsProvider(creds),
			config.WithEndpointResolverWithOptions(customResolver))
	} else {
		cfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(s.ResourceOpts.Region),
			config.WithCredentialsProvider(creds))
	}

	if err != nil {
		return nil, err
	}

	// create an S3 service client
	s3Client := s3.NewFromConfig(cfg)

	return s3Client, nil
}

func presignGetObject(client *s3.Client, bucket, objectKey string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(client, s3.WithPresignExpires(expiry))
	params := s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &objectKey,
	}
	output, err := presignClient.PresignGetObject(context.TODO(), &params)
	if err != nil {
		return "", err
	}

	return output.URL, nil
}

func presignPutObject(client *s3.Client, bucket, objectKey, ACL string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(client, s3.WithPresignExpires(expiry))
	params := s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &objectKey,
	}
	if ACL != "" {
		params.ACL = types.ObjectCannedACL(ACL)

	}
	output, err := presignClient.PresignPutObject(context.TODO(), &params)
	if err != nil {
		return "", err
	}
	return output.URL, nil
}
