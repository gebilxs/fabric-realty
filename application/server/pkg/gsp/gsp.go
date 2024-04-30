package gsp

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type OptionFunc func(*GSP)

func WithLog(log *log.Logger) OptionFunc {
	return func(g *GSP) {
		g.log = log
	}
}

type GSP struct {
	svc *session.Session // 亚马逊s3库操作器
	log *log.Logger
}

// NewGSP
// addr 确保地址包含HTTP前缀
func NewGSP(addr, access, secret, regions string, ops ...OptionFunc) (*GSP, error) {
	g := &GSP{
		log: log.Default(),
	}
	cres := credentials.NewStaticCredentials(access, secret, "")
	cfg := aws.NewConfig().WithRegion(regions).WithEndpoint(addr).WithCredentials(cres).WithS3ForcePathStyle(true)
	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}
	g.svc = sess
	for _, op := range ops {
		op(g)
	}
	return g, nil
}

func (g *GSP) PutS3ObjectWithReader(ctx context.Context, bucket, key, contentType string, data io.Reader) (string, error) {
	uploader := s3manager.NewUploader(g.svc)
	result, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		Body:        data,
	})
	if err != nil {
		return "", err
	}
	return result.Location, nil
}

func (g *GSP) PutS3Object(ctx context.Context, bucket, key, contentType string, data []byte) (string, error) {
	if len(data) < 32 {
		return "", errors.New("文件对象为空")
	}
	uploader := s3manager.NewUploader(g.svc)
	result, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		Body:        bytes.NewReader(data),
	})
	if err != nil {
		return "", err
	}
	return result.Location, nil
}
func (g *GSP) GetS3Object(ctx context.Context, bucket, key string) ([]byte, error) {
	downloader := s3manager.NewDownloader(g.svc)
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := downloader.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return buf.Bytes(), err
}

func (g *GSP) GetS3ObjectWithWriter(ctx context.Context, bucket, key string, w io.WriterAt) error {
	downloader := s3manager.NewDownloader(g.svc)
	_, err := downloader.DownloadWithContext(ctx, w, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}
