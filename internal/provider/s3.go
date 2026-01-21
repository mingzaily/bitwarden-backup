package provider

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/mingzaily/bitwarden-backup/internal/model"
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

// Backup 执行 S3 备份，返回最终存储路径
func (p *S3Provider) Backup(ctx BackupContext) (string, error) {
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
		return "", fmt.Errorf("failed to load S3 config: %w", err)
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
		return "", fmt.Errorf("failed to open source file: %w", err)
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
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// 返回 S3 路径
	return fmt.Sprintf("s3://%s/%s", dest.S3Bucket, key), nil
}

// Cleanup 清理超出保留数量的旧备份
func (p *S3Provider) Cleanup(dest model.BackupDestination, maxCount int) (int, error) {
	if maxCount <= 0 {
		return 0, nil
	}

	client, err := p.createClient(dest)
	if err != nil {
		return 0, err
	}

	// 构建前缀
	prefix := strings.TrimPrefix(dest.S3Path, "/")
	if prefix != "" {
		prefix = prefix + "/"
	}
	prefix = prefix + "backup_"

	// 列举对象
	result, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(dest.S3Bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to list objects: %w", err)
	}

	// 筛选备份文件
	var backups []types.Object
	for _, obj := range result.Contents {
		key := aws.ToString(obj.Key)
		if strings.HasSuffix(key, ".json") {
			backups = append(backups, obj)
		}
	}

	if len(backups) <= maxCount {
		return 0, nil
	}

	// 按修改时间降序排序（处理 nil 情况）
	sort.Slice(backups, func(i, j int) bool {
		if backups[i].LastModified == nil {
			return false
		}
		if backups[j].LastModified == nil {
			return true
		}
		return backups[i].LastModified.After(*backups[j].LastModified)
	})

	// 删除超出数量的旧文件
	var toDelete []types.ObjectIdentifier
	for i := maxCount; i < len(backups); i++ {
		toDelete = append(toDelete, types.ObjectIdentifier{
			Key: backups[i].Key,
		})
	}

	if len(toDelete) == 0 {
		return 0, nil
	}

	_, err = client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(dest.S3Bucket),
		Delete: &types.Delete{Objects: toDelete},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to delete objects: %w", err)
	}

	return len(toDelete), nil
}

// createClient 创建 S3 客户端
func (p *S3Provider) createClient(dest model.BackupDestination) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			dest.S3AccessKey,
			dest.S3SecretKey,
			"",
		)),
		config.WithRegion(dest.S3Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load S3 config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if dest.S3Endpoint != "" {
			o.BaseEndpoint = aws.String(dest.S3Endpoint)
			o.UsePathStyle = true
		}
	})

	return client, nil
}
