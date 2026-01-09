package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Provider S3 存储提供者
type S3Provider struct{}

// NewS3Provider 创建 S3 存储提供者
func NewS3Provider() *S3Provider {
	return &S3Provider{}
}

// Type 返回提供者类型
func (p *S3Provider) Type() string {
	return "s3"
}

// Backup 执行 S3 备份
func (p *S3Provider) Backup(ctx BackupContext) error {
	dest := ctx.Destination

	// 创建 S3 客户端配置
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			dest.S3AccessKey,
			dest.S3SecretKey,
			"",
		)),
		config.WithRegion(dest.S3Region),
	)
	if err != nil {
		return fmt.Errorf("failed to load S3 config: %w", err)
	}

	// 创建 S3 客户端
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if dest.S3Endpoint != "" {
			o.BaseEndpoint = aws.String(dest.S3Endpoint)
			o.UsePathStyle = true
		}
	})

	// 打开源文件
	file, err := os.Open(ctx.SourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	// 构建远程路径
	remotePath := strings.TrimPrefix(dest.S3Path, "/")
	if remotePath != "" {
		remotePath = remotePath + "/"
	}
	key := fmt.Sprintf("%sbackup_%s_%s.json", remotePath, ctx.TaskName, ctx.Timestamp)

	// 上传文件
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(dest.S3Bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	return nil
}

// getFileName 从路径中提取文件名
func getFileName(path string) string {
	return filepath.Base(path)
}
