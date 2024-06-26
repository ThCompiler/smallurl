// Package evjson оборачивает схему [vjson.Schema], чтобы получать вместо текстовой ошибки парсинга json,
// ошибку, заданную в качестве переменной [ErrorInvalidJSON].
package evjson

import (
	"github.com/miladibra10/vjson"
	"github.com/pkg/errors"
)

const jsonError = "could not parse json input."

var ErrorInvalidJSON = errors.New(jsonError)

type Schema struct {
	vjson.Schema
}

func NewSchema(fields ...vjson.Field) Schema {
	return Schema{vjson.NewSchema(fields...)}
}

func (s *Schema) ValidateBytes(input []byte) error {
	if err := s.Schema.ValidateBytes(input); err != nil {
		if err.Error() == jsonError {
			return ErrorInvalidJSON
		}

		return err
	}

	return nil
}
