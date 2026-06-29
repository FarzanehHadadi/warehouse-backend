package filter

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CursorData struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// ReportCursorData is used for keyset pagination on (current_quantity, product_id).
type ReportCursorData struct {
	ProductID       uint
	CurrentQuantity int
}

// StoreProductQuantityCursorData is used for keyset pagination on store product quantity rows.
type StoreProductQuantityCursorData struct {
	TotalQuantity int
	StoreID       uint
	ProductID     uint
}

func EncodeCursor(id uint, createdAt time.Time) string {
	data := fmt.Sprintf("%d|%s", id, createdAt.Format(time.RFC3339Nano))
	return base64.RawStdEncoding.EncodeToString([]byte(data))
}

func DecodeCursor(cursor string) (*CursorData, error) {
	if cursor == "" {
		return nil, nil
	}

	decoded, err := decodeCursorBytes(cursor)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(decoded), "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid cursor data")
	}

	t, err := time.Parse(time.RFC3339Nano, parts[1])
	if err != nil {
		return nil, err
	}

	var id uint
	fmt.Sscanf(parts[0], "%d", &id)
	return &CursorData{
		ID:        id,
		CreatedAt: t,
	}, nil
}

func EncodeReportCursor(productID uint, currentQuantity int) string {
	data := fmt.Sprintf("%d|%d", productID, currentQuantity)
	return base64.RawStdEncoding.EncodeToString([]byte(data))
}

func EncodeStoreProductQuantityCursor(totalQuantity int, storeID, productID uint) string {
	data := fmt.Sprintf("%d|%d|%d", totalQuantity, storeID, productID)
	return base64.RawStdEncoding.EncodeToString([]byte(data))
}

func DecodeStoreProductQuantityCursor(cursor string) (*StoreProductQuantityCursorData, error) {
	if cursor == "" {
		return nil, nil
	}

	decoded, err := decodeCursorBytes(cursor)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(decoded), "|")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid cursor data")
	}

	totalQuantity, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid cursor data")
	}

	var storeID, productID uint
	if _, err := fmt.Sscanf(parts[1], "%d", &storeID); err != nil {
		return nil, fmt.Errorf("invalid cursor data")
	}
	if _, err := fmt.Sscanf(parts[2], "%d", &productID); err != nil {
		return nil, fmt.Errorf("invalid cursor data")
	}

	return &StoreProductQuantityCursorData{
		TotalQuantity: totalQuantity,
		StoreID:       storeID,
		ProductID:     productID,
	}, nil
}

func DecodeReportCursor(cursor string) (*ReportCursorData, error) {
	if cursor == "" {
		return nil, nil
	}

	decoded, err := decodeCursorBytes(cursor)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(decoded), "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid cursor data")
	}

	var productID uint
	if _, err := fmt.Sscanf(parts[0], "%d", &productID); err != nil {
		return nil, fmt.Errorf("invalid cursor data")
	}

	currentQuantity, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid cursor data")
	}

	return &ReportCursorData{
		ProductID:       productID,
		CurrentQuantity: currentQuantity,
	}, nil
}

// BuildPaginatedResult trims an over-fetched slice and builds cursor metadata.
func BuildPaginatedResult[T any](items []T, limit int, encode func(T) string) ([]T, CursorResponse) {
	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}

	nextCursor := ""
	if hasMore && len(items) > 0 {
		nextCursor = encode(items[len(items)-1])
	}

	return items, CursorResponse{
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}
}

func decodeCursorBytes(cursor string) ([]byte, error) {
	decoded, err := base64.RawStdEncoding.DecodeString(cursor)
	if err != nil {
		decoded, err = base64.RawURLEncoding.DecodeString(cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor data")
		}
	}
	return decoded, nil
}
