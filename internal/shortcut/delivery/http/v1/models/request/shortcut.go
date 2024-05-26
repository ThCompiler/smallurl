package request

import (
	"smallurl/internal/pkg/evjson"

	"github.com/miladibra10/vjson"
)

type LongURL struct {
	// Оригинальный URL, который необходимо сократить
	OriginalURL string `json:"original_url" swaggertype:"string" example:"http://example.com"`
}

func ValidateLongURL(data []byte) error {
	schema := evjson.NewSchema(
		vjson.String("original_url").Required(),
	)

	return schema.ValidateBytes(data)
}
