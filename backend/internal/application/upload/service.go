package upload

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/blytz/live/backend/internal/infrastructure/storage/r2"
)

// Service handles file uploads
type Service struct {
	r2Client *r2.Client
}

// NewService creates a new upload service
func NewService(r2Client *r2.Client) *Service {
	return &Service{r2Client: r2Client}
}

// UploadResult represents an upload result
type UploadResult struct {
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	Key          string `json:"key"`
	ContentType  string `json:"content_type"`
	Size         int64  `json:"size"`
}

// UploadProductImage uploads a product image
func (s *Service) UploadProductImage(ctx context.Context, file io.Reader, filename string, size int64) (*UploadResult, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = ".jpg"
	}

	// Validate file type
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
	}
	if !allowedExts[ext] {
		return nil, fmt.Errorf("invalid file type: %s (allowed: jpg, png, webp, gif)", ext)
	}

	// Max 10MB
	if size > 10*1024*1024 {
		return nil, fmt.Errorf("file too large: max 10MB")
	}

	result, err := s.r2Client.Upload(ctx, file, filename, r2.UploadOptions{
		Folder:            "products",
		GenerateThumbnail: true,
		AllowedTypes:      []string{"image/"},
	})
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		URL:          result.URL,
		ThumbnailURL: result.ThumbnailURL,
		Key:          result.Key,
		ContentType:  result.ContentType,
		Size:         size,
	}, nil
}

// UploadAvatar uploads a user avatar
func (s *Service) UploadAvatar(ctx context.Context, file io.Reader, filename string, size int64) (*UploadResult, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = ".jpg"
	}

	// Validate file type
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		return nil, fmt.Errorf("invalid file type: %s", ext)
	}

	// Max 5MB
	if size > 5*1024*1024 {
		return nil, fmt.Errorf("file too large: max 5MB")
	}

	result, err := s.r2Client.Upload(ctx, file, filename, r2.UploadOptions{
		Folder:            "avatars",
		GenerateThumbnail: true,
		AllowedTypes:      []string{"image/"},
	})
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		URL:          result.URL,
		ThumbnailURL: result.ThumbnailURL,
		Key:          result.Key,
		ContentType:  result.ContentType,
		Size:         size,
	}, nil
}

// UploadStreamThumbnail uploads a stream thumbnail
func (s *Service) UploadStreamThumbnail(ctx context.Context, file io.Reader, filename string, size int64) (*UploadResult, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = ".jpg"
	}

	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		return nil, fmt.Errorf("invalid file type: %s", ext)
	}

	if size > 5*1024*1024 {
		return nil, fmt.Errorf("file too large: max 5MB")
	}

	result, err := s.r2Client.Upload(ctx, file, filename, r2.UploadOptions{
		Folder:            "streams",
		GenerateThumbnail: false,
		AllowedTypes:      []string{"image/"},
	})
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		URL:         result.URL,
		Key:         result.Key,
		ContentType: result.ContentType,
		Size:        size,
	}, nil
}

// UploadFromMultipart uploads a file from multipart form
func (s *Service) UploadFromMultipart(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (*UploadResult, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	switch folder {
	case "products":
		return s.UploadProductImage(ctx, file, fileHeader.Filename, fileHeader.Size)
	case "avatars":
		return s.UploadAvatar(ctx, file, fileHeader.Filename, fileHeader.Size)
	case "streams":
		return s.UploadStreamThumbnail(ctx, file, fileHeader.Filename, fileHeader.Size)
	default:
		return nil, fmt.Errorf("unknown folder: %s", folder)
	}
}

// DeleteFile deletes a file from storage
func (s *Service) DeleteFile(ctx context.Context, key string) error {
	return s.r2Client.Delete(ctx, key)
}

// DeleteByURL deletes a file by its public URL
func (s *Service) DeleteByURL(ctx context.Context, url string) error {
	key := s.r2Client.ExtractKey(url)
	if key == "" {
		return fmt.Errorf("invalid URL")
	}
	return s.DeleteFile(ctx, key)
}

// GetUploadPresignedURL generates a presigned URL for direct browser upload
func (s *Service) GetUploadPresignedURL(ctx context.Context, key string, expiryMinutes int) (string, error) {
	return s.r2Client.GetSignedURL(ctx, key, time.Duration(expiryMinutes)*time.Minute)
}
