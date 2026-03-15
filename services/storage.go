package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/irfanseptian/fims-backend/config"
)

// StorageService handles upload operations to Supabase Storage.
type StorageService struct {
	supabaseURL    string
	serviceRoleKey string
	bucketName     string
	httpClient     *http.Client
}

// NewStorageService creates a storage service instance from app config.
func NewStorageService(cfg *config.Config) *StorageService {
	return &StorageService{
		supabaseURL:    strings.TrimRight(strings.TrimSpace(cfg.SupabaseURL), "/"),
		serviceRoleKey: strings.TrimSpace(cfg.SupabaseServiceRoleKey),
		bucketName:     strings.TrimSpace(cfg.SupabaseStorageBucket),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsReady returns true when required Supabase Storage config is available.
func (s *StorageService) IsReady() bool {
	return s != nil && s.supabaseURL != "" && s.serviceRoleKey != "" && s.bucketName != ""
}

// UploadProfilePhoto uploads bytes to Supabase and returns public file URL.
func (s *StorageService) UploadProfilePhoto(userID string, fileName string, content []byte) (string, error) {
	if !s.IsReady() {
		return "", errors.New("Konfigurasi Supabase Storage belum lengkap")
	}

	mimeType := http.DetectContentType(content)
	if !isAllowedImageMime(mimeType) {
		return "", errors.New("Format file tidak didukung. Gunakan JPG, PNG, atau WEBP")
	}

	ext := normalizeExt(fileName, mimeType)
	objectPath := fmt.Sprintf("%s/profile-%d%s", userID, time.Now().UnixNano(), ext)
	escapedObjectPath := strings.ReplaceAll(objectPath, " ", "_")
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.supabaseURL, s.bucketName, url.PathEscape(escapedObjectPath))
	uploadURL = strings.ReplaceAll(uploadURL, "%2F", "/")

	if err := s.uploadFile(uploadURL, http.MethodPost, mimeType, content); err != nil {
		if putErr := s.uploadFile(uploadURL, http.MethodPut, mimeType, content); putErr != nil {
			return "", err
		}
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.supabaseURL, s.bucketName, escapedObjectPath)
	return publicURL, nil
}

func (s *StorageService) uploadFile(uploadURL string, method string, mimeType string, content []byte) error {
	req, err := http.NewRequestWithContext(context.Background(), method, uploadURL, bytes.NewReader(content))
	if err != nil {
		return errors.New("Gagal menyiapkan upload foto")
	}

	req.Header.Set("Authorization", "Bearer "+s.serviceRoleKey)
	req.Header.Set("apikey", s.serviceRoleKey)
	req.Header.Set("Content-Type", mimeType)
	req.Header.Set("x-upsert", "true")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return errors.New("Gagal upload foto profil")
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("Upload foto gagal (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}

func isAllowedImageMime(mimeType string) bool {
	allowed := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	return allowed[mimeType]
}

func normalizeExt(fileName, mimeType string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext != "" {
		return ext
	}

	extensions, err := mime.ExtensionsByType(mimeType)
	if err == nil && len(extensions) > 0 {
		return strings.ToLower(extensions[0])
	}

	switch mimeType {
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}
