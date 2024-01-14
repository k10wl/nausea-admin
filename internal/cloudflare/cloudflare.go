package cloudflare

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Cloudflare struct {
	client           *s3.Client
	bucketName       string
	publicStorageUrl string
}

func NewClient() *Cloudflare {
	bucketName := os.Getenv("BUCKET_NAME")
	accountId := os.Getenv("ACCOUNT_ID")
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")
	publicStorageUrl := os.Getenv("PUBLIC_STORAGE_URL")
	r2Resolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId),
				HostnameImmutable: true,
				Source:            aws.EndpointSourceCustom,
			}, nil
		},
	)
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, ""),
		),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := s3.NewFromConfig(cfg)
	return &Cloudflare{
		client:           client,
		bucketName:       bucketName,
		publicStorageUrl: publicStorageUrl,
	}
}

func (c *Cloudflare) AddObject(file io.Reader, name string) (string, error) {
	buf := make([]byte, 512) // 512 bytes are usually enough for MIME detection
	n, err := io.ReadFull(file, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", err
	}

	contentType := http.DetectContentType(buf[:n])

	var reader io.Reader
	if s, ok := file.(io.Seeker); ok {
		_, err = s.Seek(0, io.SeekStart)
		if err != nil {
			return "", err
		}
		reader = file
	} else {
		reader = io.MultiReader(bytes.NewReader(buf[:n]), file)
	}

	key := aws.String(name)

	_, err = c.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         key,
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	url := fmt.Sprintf("%s/%s", c.publicStorageUrl, *key)

	return url, err
}

func (c *Cloudflare) RemoveObject(name string) error {
	_, err := c.client.DeleteObject(
		context.TODO(),
		&s3.DeleteObjectInput{Key: &name, Bucket: &c.bucketName},
	)
	return err
}

func (c *Cloudflare) ParseURLKey(url string) string {
	return strings.TrimPrefix(url, fmt.Sprintf("%s/", c.publicStorageUrl))
}
