package handlers

import (
	"net/http"
	"smallurl/internal/shortcut/delivery/http/tools"
	"smallurl/internal/shortcut/delivery/http/v1/models/response"
	md "smallurl/internal/shortcut/delivery/mocks"
	sr "smallurl/internal/shortcut/repository"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"github.com/steinfletcher/apitest"
	"go.uber.org/mock/gomock"
)

var testError = errors.New("test error")

type ApiShortcutSuite struct {
	suite.Suite
	handlers       *ShortcutHandlers
	mockController *gomock.Controller
	usc            *md.MockUsecase
	router         *gin.Engine
}

func (ss *ApiShortcutSuite) BeforeEach(t provider.T) {
	ss.mockController = gomock.NewController(t)
	ss.usc = md.NewMockUsecase(ss.mockController)
	ss.handlers = NewShortcutHandlers(ss.usc)
	ss.router = gin.New()
	gin.SetMode(gin.ReleaseMode)
}

func (ss *ApiShortcutSuite) AfterEach(t provider.T) {
	ss.mockController.Finish()
}

func (ss *ApiShortcutSuite) TestGetLong(t provider.T) {
	t.Title("Тестирование апи метода GetLong: GET /:short_url")
	const path = "/api/v1/:" + ShortURLParam

	ss.router.GET(path, ss.handlers.GetLong)

	t.Run("Успешное получение длинного URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		ss.usc.EXPECT().GetLongURL(shortURL).Times(1).Return(longURL, nil)

		t.NewStep("Тестирование")
		resp := apitest.New().
			Handler(ss.router).
			Getf("/api/v1/%s", shortURL).
			Expect(t).
			Status(http.StatusSeeOther).
			End()

		t.NewStep("Проверка результатов")
		t.Require().Contains(resp.Response.Header["Location"], longURL)
	})

	t.Run("Получение длинного URL с ошибкой", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		ss.usc.EXPECT().GetLongURL(shortURL).Times(1).Return("", testError)

		t.NewStep("Тестирование")
		resp := apitest.New().
			Handler(ss.router).
			Getf("/api/v1/%s", shortURL).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()

		t.NewStep("Проверка результатов")
		var err tools.Error

		resp.JSON(&err)
		t.Require().NotEmpty(err.Error)
	})

	t.Run("Получение длинного URL с не существующему короткому", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		ss.usc.EXPECT().GetLongURL(shortURL).Times(1).Return("", sr.ErrorURLNotFound)

		t.NewStep("Тестирование")
		apitest.New().
			Handler(ss.router).
			Getf("/api/v1/%s", shortURL).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})
}

func (ss *ApiShortcutSuite) TestGenShort(t provider.T) {
	t.Title("Тестирование апи метода GenShort: POST /shorten")
	const path = "/api/v1/shorten"

	ss.router.POST(path, ss.handlers.GenShort)

	t.Run("Успешное создания короткого URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		ss.usc.EXPECT().GetShortURL(longURL).Times(1).Return(shortURL, nil)

		t.NewStep("Тестирование")
		resp := apitest.New().
			Handler(ss.router).
			Post(path).
			Bodyf("{ \"original_url\": \"%s\" }", longURL).
			Expect(t).
			Status(http.StatusCreated).
			End()

		t.NewStep("Проверка результата")
		var short response.Result
		resp.JSON(&short)

		t.Require().Equal(shortURL, short.ShortURL)
	})

	t.Run("Создание короткого URL с неверным телом запроса", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"

		t.NewStep("Тестирование с отсутствием тела запроса")
		apitest.New().
			Handler(ss.router).
			Post(path).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

		t.NewStep("Тестирование с неверным полем тела запроса")
		apitest.New().
			Handler(ss.router).
			Post(path).
			Bodyf("{ \"long_url\": \"%s\" }", longURL).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Создание короткого URL с некорректным длинным URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "httpexample.com"

		t.NewStep("Тестирование")
		apitest.New().
			Handler(ss.router).
			Post(path).
			Bodyf("{ \"original_url\": \"%s\" }", longURL).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Создание короткого URL с ошибкой", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"

		t.NewStep("Инициализация мока")
		ss.usc.EXPECT().GetShortURL(longURL).Times(1).Return("", testError)

		t.NewStep("Тестирование")
		resp := apitest.New().
			Handler(ss.router).
			Post(path).
			Bodyf("{ \"original_url\": \"%s\" }", longURL).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()

		t.NewStep("Проверка результатов")
		var err tools.Error

		resp.JSON(&err)
		t.Require().NotEmpty(err.Error)
	})
}

func TestRunApiShortcutTest(t *testing.T) {
	suite.RunSuite(t, new(ApiShortcutSuite))
}
