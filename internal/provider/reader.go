package provider

import (
	"context"
)

// DataReader is the common contract
type DataReader[T any] interface {
	Read(ctx context.Context) ([]*T, error)
}
