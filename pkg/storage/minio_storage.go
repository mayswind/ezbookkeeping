package storage

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// MinIOObjectStorage represents MinIO object storage
type MinIOObjectStorage struct {
	minIOClient *minio.Client
	minIOConfig *settings.MinIOConfig
	rootPath    string
}

// NewMinIOObjectStorage returns a MinIO object storage
func NewMinIOObjectStorage(config *settings.Config, pathPrefix string) (*MinIOObjectStorage, error) {
	minIOConfig := config.MinIOConfig

	minIOClient, err := minio.New(minIOConfig.Endpoint, &minio.Options{
		Region:    minIOConfig.Location,
		Creds:     credentials.NewStaticV4(minIOConfig.AccessKeyID, minIOConfig.SecretAccessKey, ""),
		Secure:    minIOConfig.UseSSL,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: minIOConfig.SkipTLSVerify}},
	})

	if err != nil {
		return nil, err
	}

	storage := &MinIOObjectStorage{
		minIOClient: minIOClient,
		minIOConfig: minIOConfig,
		rootPath:    minIOConfig.RootPath,
	}

	storage.rootPath = storage.getFinalPath(pathPrefix)
	storage.rootPath = strings.ReplaceAll(storage.rootPath, "\\", "/")

	ctx := context.Background()
	exists, err := minIOClient.BucketExists(ctx, minIOConfig.Bucket)

	if err != nil {
		return nil, err
	}

	if !exists {
		err := minIOClient.MakeBucket(ctx, minIOConfig.Bucket, minio.MakeBucketOptions{
			Region: minIOConfig.Location,
		})

		if err != nil {
			return nil, err
		}
	}

	return storage, nil
}

// Exists returns whether the file exists
func (s *MinIOObjectStorage) Exists(path string) (bool, error) {
	ctx := context.Background()
	objectInfo, err := s.minIOClient.StatObject(ctx, s.minIOConfig.Bucket, s.getFinalPath(path), minio.StatObjectOptions{})

	if err == nil && !objectInfo.IsDeleteMarker {
		return true, nil
	}

	return false, err
}

// Read returns the object instance according to specified the file path
func (s *MinIOObjectStorage) Read(path string) (ObjectInStorage, error) {
	ctx := context.Background()
	return s.minIOClient.GetObject(ctx, s.minIOConfig.Bucket, s.getFinalPath(path), minio.GetObjectOptions{})
}

// Save returns whether save the object instance successfully
func (s *MinIOObjectStorage) Save(path string, object ObjectInStorage) error {
	ctx := context.Background()
	_, err := s.minIOClient.PutObject(ctx, s.minIOConfig.Bucket, s.getFinalPath(path), object, -1, minio.PutObjectOptions{})

	return err
}

// Delete returns whether delete the object according to specified the file path successfully
func (s *MinIOObjectStorage) Delete(path string) error {
	ctx := context.Background()
	return s.minIOClient.RemoveObject(ctx, s.minIOConfig.Bucket, s.getFinalPath(path), minio.RemoveObjectOptions{})
}

func (s *MinIOObjectStorage) getFinalPath(path string) string {
	rootPath := s.rootPath

	if len(rootPath) < 1 || rootPath[len(rootPath)-1] != '/' {
		rootPath = rootPath + "/"
	}

	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	path = strings.ReplaceAll(path, "\\", "/")

	return rootPath + path
}
