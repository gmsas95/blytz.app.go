package r2

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

// Client represents an R2 storage client
type Client struct {
	s3Client  *s3.Client
	bucket    string
	publicURL string
	cdnURL    string
}

// Config holds R2 configuration
type Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	PublicURL       string // e.g., https://pub-xxx.r2.dev
	CDNURL          string // e.g., https://cdn.blytz.app
}

// NewClient creates a new R2 client
func NewClient(cfg Config) (*Client, error) {
	// Create custom resolver for R2 endpoint
	endpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID),
		}, nil
	})

	awsCfg := aws.Config{
		EndpointResolverWithOptions: endpointResolver,
		Credentials:                 credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Region:                      "auto",
	}

	s3Client := s3.NewFromConfig(awsCfg)

	return &Client{
		s3Client:  s3Client,
		bucket:    cfg.BucketName,
		publicURL: cfg.PublicURL,
		cdnURL:    cfg.CDNURL,
	}, nil
}

// UploadResult represents the result of an upload
type UploadResult struct {
	URL          string
	ThumbnailURL string
	Key          string
	ContentType  string
	Size         int64
}

// UploadOptions represents upload options
type UploadOptions struct {
	ContentType       string
	Folder            string // e.g., "products", "avatars", "streams"
	GenerateThumbnail bool
	MaxSize           int64
	AllowedTypes      []string
}

// Upload uploads a file to R2
func (c *Client) Upload(ctx context.Context, reader io.Reader, filename string, opts UploadOptions) (*UploadResult, error) {
	// Validate content type
	contentType := opts.ContentType
	if contentType == "" {
		contentType = mime.TypeByExtension(path.Ext(filename))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	// Check allowed types
	if len(opts.AllowedTypes) > 0 {
		allowed := false
		for _, t := range opts.AllowedTypes {
			if strings.HasPrefix(contentType, t) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("content type %s not allowed", contentType)
		}
	}

	// Generate unique key
	ext := path.Ext(filename)
	if ext == "" {
		ext = ".bin"
	}
	key := fmt.Sprintf("%s/%s%s", opts.Folder, uuid.New().String(), ext)

	// Upload to R2
	putInput := &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
		Metadata: map[string]string{
			"original-filename": filename,
			"uploaded-at":       time.Now().UTC().Format(time.RFC3339),
		},
	}

	_, err := c.s3Client.PutObject(ctx, putInput)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to R2: %w", err)
	}

	// Generate URLs
	result := &UploadResult{
		URL:         c.getPublicURL(key),
		Key:         key,
		ContentType: contentType,
	}

	// Generate thumbnail URL for images
	if opts.GenerateThumbnail && strings.HasPrefix(contentType, "image/") {
		result.ThumbnailURL = c.getThumbnailURL(key)
	}

	return result, nil
}

// UploadBytes uploads bytes to R2
func (c *Client) UploadBytes(ctx context.Context, data []byte, filename string, opts UploadOptions) (*UploadResult, error) {
	return c.Upload(ctx, strings.NewReader(string(data)), filename, opts)
}

// Delete deletes a file from R2
func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from R2: %w", err)
	}
	return nil
}

// GetSignedURL generates a presigned URL for temporary access
func (c *Client) GetSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.s3Client)

	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiry))

	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return req.URL, nil
}

// Copy copies an object within R2
func (c *Client) Copy(ctx context.Context, sourceKey, destKey string) error {
	copySource := fmt.Sprintf("%s/%s", c.bucket, url.PathEscape(sourceKey))

	_, err := c.s3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(c.bucket),
		CopySource: aws.String(copySource),
		Key:        aws.String(destKey),
	})

	if err != nil {
		return fmt.Errorf("failed to copy object: %w", err)
	}

	return nil
}

// ListObjects lists objects in a prefix
func (c *Client) ListObjects(ctx context.Context, prefix string, maxKeys int32) ([]types.Object, error) {
	result, err := c.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(c.bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: maxKeys,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	return result.Contents, nil
}

// GetObject retrieves an object
func (c *Client) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return result.Body, nil
}

// Helper methods

func (c *Client) getPublicURL(key string) string {
	if c.cdnURL != "" {
		return fmt.Sprintf("%s/%s", c.cdnURL, key)
	}
	return fmt.Sprintf("%s/%s", c.publicURL, key)
}

func (c *Client) getThumbnailURL(key string) string {
	// For now, return the same URL
	// In production, you might use Cloudflare Images or Image Resizing
	return c.getPublicURL(key)
}

// ExtractKey extracts the key from a public URL
func (c *Client) ExtractKey(publicURL string) string {
	// Remove base URL to get key
	if strings.HasPrefix(publicURL, c.cdnURL) {
		return strings.TrimPrefix(publicURL, c.cdnURL+"/")
	}
	if strings.HasPrefix(publicURL, c.publicURL) {
		return strings.TrimPrefix(publicURL, c.publicURL+"/")
	}
	return publicURL
}
