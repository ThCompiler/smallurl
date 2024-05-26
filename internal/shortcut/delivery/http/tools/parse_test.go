package tools

import (
	"io"
	"net/http"
	"smallurl/internal/pkg/evjson"
	"smallurl/pkg/logger"
	"strings"
	"testing"

	"github.com/miladibra10/vjson"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/pkg/errors"
)

var testError = errors.New("test error")

type request struct {
	Input string `json:"input"`
}

func ValidateInput(data []byte) error {
	schema := evjson.NewSchema(
		vjson.String("input").Required(),
	)

	return schema.ValidateBytes(data)
}

const (
	invalidJson      = "{ asdas"
	unknownFieldJson = "{ \"field\" : 1 }"
)

func TestParse(t *testing.T) {
	runner.Run(t, "Тестирование ParseRequestBody", func(t provider.T) {
		t.NewStep("Инициализация исходных данных")
		body := "{ \"input\": \"input\"}"
		incorrectBody := "{ \"input\": 1 }"

		t.WithNewStep("Успешное получение тела запроса", func(t provider.StepCtx) {
			t.NewStep("Инициализация тела")
			b := io.NopCloser(strings.NewReader(body))

			t.NewStep("Тестирование")
			var input request
			code, err := ParseRequestBody(b, &input, ValidateInput, logger.DefaultLogger)

			t.NewStep("Проверка результата")
			t.Require().NoError(err)
			t.Require().Equal(http.StatusOK, code)
		})

		t.WithNewStep("Получение тела запроса с ошибкой", func(t provider.StepCtx) {
			t.NewStep("Тестирование")
			var input request
			code, err := ParseRequestBody(errReader(1), &input, ValidateInput, logger.DefaultLogger)

			t.NewStep("Проверка результата")
			t.Require().ErrorIs(err, ErrorCannotReadBody)
			t.Require().Equal(http.StatusInternalServerError, code)
		})

		t.WithNewStep("Получение тела запроса с некорректным JSON", func(t provider.StepCtx) {
			t.NewStep("Инициализация тела")
			b := io.NopCloser(strings.NewReader(invalidJson))

			t.NewStep("Тестирование")
			var input request
			code, err := ParseRequestBody(b, &input, ValidateInput, logger.DefaultLogger)

			t.NewStep("Проверка результата")
			t.Require().ErrorIs(err, ErrorIncorrectBodyContent)
			t.Require().Equal(http.StatusBadRequest, code)
		})

		t.WithNewStep("Получение тела запроса с неверными JSON-полями", func(t provider.StepCtx) {
			t.NewStep("Инициализация тела")
			b := io.NopCloser(strings.NewReader(unknownFieldJson))

			t.NewStep("Тестирование")
			var input request
			code, err := ParseRequestBody(b, &input, ValidateInput, logger.DefaultLogger)

			t.NewStep("Проверка результата")
			t.Require().Error(err)
			t.Require().Equal(http.StatusBadRequest, code)
		})

		t.WithNewStep("Получение тела запроса с ошибкой парсинга JSON", func(t provider.StepCtx) {
			t.NewStep("Инициализация тела")
			b := io.NopCloser(strings.NewReader(incorrectBody))

			t.NewStep("Тестирование")
			var input request
			code, err := ParseRequestBody(b, &input, func([]byte) error { return nil }, logger.DefaultLogger)

			t.NewStep("Проверка результата")
			t.Require().ErrorIs(err, ErrorIncorrectBodyContent)
			t.Require().Equal(http.StatusBadRequest, code)
		})
	})
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, testError
}

func (errReader) Close() error {
	return testError
}
