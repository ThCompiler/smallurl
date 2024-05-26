package v1

import (
	"context"
	md "smallurl/internal/shortcut/delivery/mocks"
	sr "smallurl/internal/shortcut/repository"
	generatedv1 "smallurl/pkg/grpc/v1"
	"smallurl/pkg/logger"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var testError = errors.New("test error")

type GrpcShortcutSuite struct {
	suite.Suite
	service        *ShortcutService
	mockController *gomock.Controller
	usc            *md.MockUsecase
}

func (gss *GrpcShortcutSuite) BeforeEach(t provider.T) {
	gss.mockController = gomock.NewController(t)
	gss.usc = md.NewMockUsecase(gss.mockController)
	gss.service = NewShortcutService(gss.usc, logger.DefaultLogger)
}

func (gss *GrpcShortcutSuite) AfterEach(t provider.T) {
	gss.mockController.Finish()
}

func (gss *GrpcShortcutSuite) TestGetLongURL(t provider.T) {
	t.Title("Тестирование grpc метода GetLongURL")

	t.Run("Успешное получение длинного URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		gss.usc.EXPECT().GetLongURL(shortURL).Times(1).Return(longURL, nil)

		t.NewStep("Тестирование")
		res, err := gss.service.GetLongURL(context.Background(), &generatedv1.ShortUrl{ShortUrl: shortURL})

		t.NewStep("Проверка результатов")
		t.Require().NoError(err)
		t.Require().Equal(longURL, res.LongUrl)
	})

	t.Run("Получение длинного URL с ошибкой", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		gss.usc.EXPECT().GetLongURL(shortURL).Times(1).Return("", testError)

		t.NewStep("Тестирование")
		_, err := gss.service.GetLongURL(context.Background(), &generatedv1.ShortUrl{ShortUrl: shortURL})

		t.NewStep("Проверка результатов")
		t.Require().ErrorIs(err, status.Error(codes.Internal, testError.Error()))
	})

	t.Run("Получение длинного URL с не существующему короткому", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		gss.usc.EXPECT().GetLongURL(shortURL).Times(1).Return("", sr.ErrorURLNotFound)

		t.NewStep("Тестирование")
		_, err := gss.service.GetLongURL(context.Background(), &generatedv1.ShortUrl{ShortUrl: shortURL})

		t.NewStep("Проверка результатов")
		t.Require().ErrorIs(err, status.Error(codes.NotFound, sr.ErrorURLNotFound.Error()))
	})
}

func (gss *GrpcShortcutSuite) TestGetShortURL(t provider.T) {
	t.Title("Тестирование grpc метода GetShortURL")

	t.Run("Успешное создания короткого URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		gss.usc.EXPECT().GetShortURL(longURL).Times(1).Return(shortURL, nil)

		t.NewStep("Тестирование")
		res, err := gss.service.GetShortURL(context.Background(), &generatedv1.LongUrl{LongUrl: longURL})

		t.NewStep("Проверка результатов")
		t.Require().NoError(err)
		t.Require().Equal(shortURL, res.ShortUrl)
	})

	t.Run("Создание короткого URL с ошибкой", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"

		t.NewStep("Инициализация мока")
		gss.usc.EXPECT().GetShortURL(longURL).Times(1).Return("", testError)

		t.NewStep("Тестирование")
		_, err := gss.service.GetShortURL(context.Background(), &generatedv1.LongUrl{LongUrl: longURL})

		t.NewStep("Проверка результатов")
		t.Require().ErrorIs(err, status.Error(codes.Internal, testError.Error()))
	})

	t.Run("Создание короткого URL с некорректным длинным URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "httpexample.com"

		t.NewStep("Тестирование")
		_, err := gss.service.GetShortURL(context.Background(), &generatedv1.LongUrl{LongUrl: longURL})

		t.NewStep("Проверка результатов")
		t.Require().ErrorIs(err, status.Error(codes.InvalidArgument, ErrorURLNotValid.Error()))
	})
}

func TestRunGrpcShortcutSuiteTest(t *testing.T) {
	suite.RunSuite(t, new(GrpcShortcutSuite))
}
