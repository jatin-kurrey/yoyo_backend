package services

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"yoyo-server/internal/config"
	"yoyo-server/internal/models"
	"yoyo-server/internal/repositories"

	"github.com/google/uuid"
	_ "golang.org/x/image/webp"
)

const maxImageDimension = 8000

type UploadService struct {
	cfg      *config.Config
	repo     *repositories.MediaAssetRepository
	audit    *AuditService
	provider StorageProvider
}

func NewUploadService(cfg *config.Config, repo *repositories.MediaAssetRepository, audit *AuditService, provider StorageProvider) *UploadService {
	return &UploadService{
		cfg:      cfg,
		repo:     repo,
		audit:    audit,
		provider: provider,
	}
}

func (s *UploadService) Save(ctx context.Context, adminID uuid.UUID, fileHeader *multipart.FileHeader, folder string, ip string) (*models.MediaAsset, error) {
	// 1. Basic checks
	if fileHeader.Size > s.cfg.MaxUploadSizeBytes {
		return nil, fmt.Errorf("file size exceeds limit of %d bytes", s.cfg.MaxUploadSizeBytes)
	}

	// 2. Open file
	source, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer source.Close()

	// 3. Detect MIME type accurately
	buffer := make([]byte, 512)
	if _, err := source.Read(buffer); err != nil && err != io.EOF {
		return nil, err
	}
	if _, err := source.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	mimeType := http.DetectContentType(buffer)
	if !allowedMimeType(mimeType) {
		return nil, fmt.Errorf("unsupported file type: %s", mimeType)
	}

	// 4. Validate image dimensions
	if strings.HasPrefix(mimeType, "image/") && mimeType != "image/svg+xml" {
		if _, err := source.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		imgCfg, _, err := image.DecodeConfig(source)
		if err == nil {
			if imgCfg.Width > maxImageDimension || imgCfg.Height > maxImageDimension {
				return nil, fmt.Errorf("image dimensions (%dx%d) exceed maximum allowed (%dx%d)",
					imgCfg.Width, imgCfg.Height, maxImageDimension, maxImageDimension)
			}
		}
		if _, err := source.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
	}

	// 5. Sanitize folder path (prevent path traversal)
	if folder != "" {
		folder = strings.TrimSpace(folder)
		folder = strings.Trim(folder, "/")
		if strings.Contains(folder, "..") || strings.HasPrefix(folder, "/") {
			return nil, fmt.Errorf("invalid folder path: %s", folder)
		}
	}

	// 6. Generate unique filename
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		if strings.Contains(mimeType, "svg") {
			ext = ".svg"
		} else if strings.Contains(mimeType, "webp") {
			ext = ".webp"
		} else if strings.Contains(mimeType, "png") {
			ext = ".png"
		} else {
			ext = ".jpg"
		}
	}
	
	fileName := fmt.Sprintf("%s%s", uuid.NewString(), ext)
	key := fileName
	if folder != "" {
		key = fmt.Sprintf("%s/%s", folder, fileName)
	}

	// 7. Upload to provider
	url, err := s.provider.Save(ctx, key, source, mimeType)
	if err != nil {
		return nil, err
	}

	// 8. Sanitize original filename for storage
	safeOriginalName := sanitizeFilename(fileHeader.Filename)

	// 9. Save record to DB
	asset := &models.MediaAsset{
		URL:              url,
		StorageKey:       key,
		Filename:         fileName,
		OriginalFilename: safeOriginalName,
		MimeType:         mimeType,
		SizeBytes:        fileHeader.Size,
		StorageProvider:  s.provider.Name(),
		UploadedByID:     adminID,
		Folder:           folder,
	}

	if err := s.repo.Create(ctx, asset); err != nil {
		_ = s.provider.Delete(ctx, key)
		return nil, err
	}

	// Audit Log
	s.audit.Log(ctx, &adminID, "media.upload", "media", map[string]interface{}{
		"id":     asset.ID,
		"url":    asset.URL,
		"folder": asset.Folder,
		"name":   asset.OriginalFilename,
	}, ip)

	return asset, nil
}

func (s *UploadService) List(ctx context.Context, page, limit int) ([]models.MediaAsset, int64, error) {
	return s.repo.List(ctx, page, limit)
}

func (s *UploadService) Delete(ctx context.Context, adminID uuid.UUID, id uuid.UUID, ip string) error {
	asset, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete from storage
	if err := s.provider.Delete(ctx, asset.StorageKey); err != nil {
		log.Printf("Warning: Failed to delete file from storage: %v\n", err)
	}

	// Delete from DB
	if err := s.repo.Delete(ctx, asset); err != nil {
		return err
	}

	// Audit Log
	s.audit.Log(ctx, &adminID, "media.delete", "media", map[string]interface{}{
		"id":   asset.ID,
		"name": asset.OriginalFilename,
	}, ip)

	return nil
}

func allowedMimeType(mime string) bool {
	valid := []string{
		"image/jpeg",
		"image/png",
		"image/webp",
		"image/svg+xml",
		"image/gif",
	}
	for _, v := range valid {
		if strings.HasPrefix(mime, v) {
			return true
		}
	}
	return false
}

func sanitizeFilename(name string) string {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	cleaned := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == ' ' || r == '-' || r == '_' || r == '.' {
			return r
		}
		return '_'
	}, base)
	cleaned = strings.TrimSpace(cleaned)
	if cleaned == "" || cleaned == "." {
		return "unnamed"
	}
	return cleaned + strings.ToLower(ext)
}
