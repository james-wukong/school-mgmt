// internal/readers/factory.go
package provider

import (
	"fmt"

	"gorm.io/gorm"
)

type SourceType string

const (
	SourceText SourceType = "text"
	SourceDB   SourceType = "db"
	SourceAPI  SourceType = "api"
	SourceCSV  SourceType = "csv"
)

type ReaderConfig struct {
	Source        SourceType
	DB            *gorm.DB // used when Source == SourceDB
	FilePath      string   // used when Source == SourceCSV
	BaseURL       string   // used when Source == SourceAPI
	Text          string   // used when Source == SourceText
	CSVSkipHeader bool
}

func NewReader[T any](cfg ReaderConfig) (DataReader[T], error) {
	switch cfg.Source {
	case SourceText:
		return NewTextReader[T](cfg.Text), nil
	case SourceCSV:
		return NewCSVReader[T](cfg.FilePath, cfg.CSVSkipHeader), nil
	// case SourceAPI:
	// 	return NewAPIReader[T](cfg.BaseURL), nil
	// case SourceDB:
	// 	return NewDBReader[T](cfg.DB), nil
	default:
		return nil, fmt.Errorf("unknown source type: %s", cfg.Source)
	}
}
