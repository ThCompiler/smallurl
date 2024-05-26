// Package tools содержит функционал для обработки запросов, используемый хэндлерами веб-сервиса.
package tools

import (
	"encoding/json"
	"io"
	"net/http"
	"smallurl/internal/pkg/evjson"
	"smallurl/pkg/logger"

	"github.com/pkg/errors"
)

var (
	ErrorCannotReadBody       = errors.New("can't read body")
	ErrorIncorrectBodyContent = errors.New("incorrect body content")
)

func ParseRequestBody(reqBody io.ReadCloser, out any, validation func([]byte) error, l logger.Interface) (int, error) {
	body, err := io.ReadAll(reqBody)
	if err != nil {
		l.Error("[HTTP} - %s", errors.Wrapf(err, "can't read body"))

		return http.StatusInternalServerError, ErrorCannotReadBody
	}

	// Проверка корректности тела запроса
	if err := validation(body); err != nil {
		if errors.Is(err, evjson.ErrorInvalidJSON) {
			l.Warn("[HTTP} - %s", errors.Wrapf(err, "try parse body json"))

			return http.StatusBadRequest, ErrorIncorrectBodyContent
		}

		return http.StatusBadRequest, errors.Wrapf(err, "in body error")
	}

	// Получение значения тела запроса
	if err := json.Unmarshal(body, out); err != nil {
		l.Warn("[HTTP} - %s", errors.Wrapf(err, "try parse create request entity"))

		return http.StatusBadRequest, ErrorIncorrectBodyContent
	}

	return http.StatusOK, nil
}
