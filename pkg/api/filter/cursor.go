package filter

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

type CursorData struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func EncodeCursor(id uint, createdAt time.Time) string {
	data := fmt.Sprintf("%d|%s", id, createdAt)
	return base64.RawStdEncoding.EncodeToString([]byte(data))

}

func DecodeCursor(cursor string) (*CursorData, error) {
	if cursor == "" {
		return nil, nil
	}

	decoded, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(string(decoded), "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid cursor data")
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
