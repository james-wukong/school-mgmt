package provider

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type FileValidationConfig struct {
	AllowedExtensions []string // e.g. []string{".csv", ".tsv"}
	MaxSizeBytes      int64    // e.g. 5 * 1024 * 1024 for 5MB
}

var (
	ErrInvalidExtension = errors.New("invalid file extension")
	ErrFileTooLarge     = errors.New("file too large")
)

func ValidateFile(filePath string, cfg FileValidationConfig) error {
	// ── extension check ───────────────────────────────────────────────────────
	if len(cfg.AllowedExtensions) > 0 {
		ext := strings.ToLower(filepath.Ext(filePath))
		allowed := false
		if slices.Contains(cfg.AllowedExtensions, ext) {
			allowed = true
		}
		if !allowed {
			return fmt.Errorf("%w: got %q, allowed: %v",
				ErrInvalidExtension, ext, cfg.AllowedExtensions)
		}
	}

	// ── size check ────────────────────────────────────────────────────────────
	if cfg.MaxSizeBytes > 0 {
		info, err := os.Stat(filePath)
		if err != nil {
			return fmt.Errorf("cannot stat file: %w", err)
		}
		if info.Size() > cfg.MaxSizeBytes {
			return fmt.Errorf("%w: got %s, max allowed: %s",
				ErrFileTooLarge,
				formatBytes(info.Size()),
				formatBytes(cfg.MaxSizeBytes),
			)
		}
	}

	return nil
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
