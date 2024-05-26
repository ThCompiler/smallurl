package usecase

import (
	"fmt"
	sr "smallurl/internal/shortcut/repository"
	mu "smallurl/internal/shortcut/usecase/mocks"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"go.uber.org/mock/gomock"
)

var testError = errors.New("test error")

type UsecaseSuite struct {
	suite.Suite
	usc            *Usecase
	mockController *gomock.Controller
	rep            *mu.MockRepository
	shortcut       *mu.MockShortcut
}

func (us *UsecaseSuite) BeforeEach(t provider.T) {
	us.mockController = gomock.NewController(t)
	us.rep = mu.NewMockRepository(us.mockController)
	us.shortcut = mu.NewMockShortcut(us.mockController)
	us.usc = NewUsecase(us.rep, func() Shortcut { return us.shortcut })
}

func (us *UsecaseSuite) AfterEach(t provider.T) {
	us.mockController.Finish()
}

func (us *UsecaseSuite) TestGetLongURL(t provider.T) {
	t.Title("Тестирование метода GetLongURL")

	t.Run("Успешное получение длинного URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		us.rep.EXPECT().GetLongURL(shortURL).Times(1).Return(longURL, nil)

		t.NewStep("Тестирование")
		res, err := us.usc.GetLongURL(shortURL)

		t.NewStep("Проверка результатов")
		t.Require().NoError(err)
		t.Require().Equal(longURL, res)
	})

	t.Run("Получение длинного URL с ошибкой", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		us.rep.EXPECT().GetLongURL(shortURL).Times(1).Return("", testError)

		t.NewStep("Тестирование")
		_, err := us.usc.GetLongURL(shortURL)

		t.NewStep("Проверка результатов")
		t.Require().ErrorIs(err, testError)
	})
}

func (us *UsecaseSuite) TestGetShortURL(t provider.T) {
	t.Title("Тестирование метода GetShortURL")

	t.Run("Успешное добавление длинного URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		us.shortcut.EXPECT().GetShort(longURL, LenShortURL).Times(1).Return(shortURL)
		us.rep.EXPECT().GetLongURL(shortURL).Times(1).Return("", sr.ErrorURLNotFound)
		us.rep.EXPECT().AddURL(shortURL, longURL).Times(1).Return(nil)

		t.NewStep("Тестирование")
		res, err := us.usc.GetShortURL(longURL)

		t.NewStep("Проверка результатов")
		t.Require().NoError(err)
		t.Require().Equal(shortURL, res)
	})

	t.Run("Добавление длинного URL, уже существующего в хранилище", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		us.shortcut.EXPECT().GetShort(longURL, LenShortURL).Times(1).Return(shortURL)
		us.rep.EXPECT().GetLongURL(shortURL).Times(1).Return(longURL, nil)

		t.NewStep("Тестирование")
		res, err := us.usc.GetShortURL(longURL)

		t.NewStep("Проверка результатов")
		t.Require().NoError(err)
		t.Require().Equal(shortURL, res)
	})

	t.Run("Добавление длинного URL c ошибкой добавления в хранилище", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"
		newShortUrl := "asdsvasdaed"

		t.NewStep("Инициализация мока")
		us.shortcut.EXPECT().GetShort(longURL, LenShortURL).Times(1).Return(shortURL)
		us.rep.EXPECT().GetLongURL(shortURL).Times(1).Return("other", nil)
		us.shortcut.EXPECT().GetShort(shortURL, LenShortURL).Times(1).Return(newShortUrl)
		us.rep.EXPECT().AddURL(newShortUrl, longURL).Times(1).Return(testError)

		t.NewStep("Тестирование")
		_, err := us.usc.GetShortURL(longURL)

		t.NewStep("Проверка результатов")
		t.Require().ErrorIs(err, testError)
	})

	t.Run("Добавление длинного URL, для которого оказалась свободна вторая генерация короткого URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"
		newShortUrl := "asdsvasdaed"

		t.NewStep("Инициализация мока")
		us.shortcut.EXPECT().GetShort(longURL, LenShortURL).Times(1).Return(shortURL)
		us.rep.EXPECT().GetLongURL(shortURL).Times(1).Return("other", nil)
		us.shortcut.EXPECT().GetShort(shortURL, LenShortURL).Times(1).Return(newShortUrl)
		us.rep.EXPECT().AddURL(newShortUrl, longURL).Times(1).Return(nil)

		t.NewStep("Тестирование")
		res, err := us.usc.GetShortURL(longURL)

		t.NewStep("Проверка результатов")
		t.Require().NoError(err)
		t.Require().Equal(newShortUrl, res)
	})

	t.Run("Добавление длинного URL, для которого оказалась свободна вторая генерация "+
		"короткого URL из-за параллельности запросов", func(t provider.T) {
		t.Description("Тест рассматривает случай, когда метод получения длинного URL на основе короткого URL " +
			"сообщил, что данный короткий не занят, но когда метод дошёл до добавления этой пары URL, другой поток" +
			"уже добавил данный короткий URL")
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"
		newShortUrl := "asdsvasdaed"

		t.NewStep("Инициализация мока")
		us.shortcut.EXPECT().GetShort(longURL, LenShortURL).Times(1).Return(shortURL)
		us.rep.EXPECT().GetLongURL(shortURL).Times(1).Return("", sr.ErrorURLNotFound)
		us.rep.EXPECT().AddURL(shortURL, longURL).Times(1).Return(sr.ErrorURLAlreadyExists)
		us.shortcut.EXPECT().GetShort(shortURL, LenShortURL).Times(1).Return(newShortUrl)
		us.rep.EXPECT().AddURL(newShortUrl, longURL).Times(1).Return(nil)

		t.NewStep("Тестирование")
		res, err := us.usc.GetShortURL(longURL)

		t.NewStep("Проверка результатов")
		t.Require().NoError(err)
		t.Require().Equal(newShortUrl, res)
	})

	t.Run("Добавление длинного URL, для которого оказалась свободна 10 генерация короткого URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		firstShortURL := "asdsvasdasd"
		lastShortUrl := "aedsvasdaed"

		t.NewStep("Инициализация мока")
		us.shortcut.EXPECT().GetShort(longURL, LenShortURL).Times(1).Return(firstShortURL)
		us.rep.EXPECT().GetLongURL(firstShortURL).Times(1).Return("other", nil)

		for i := 0; i < 9; i++ {
			if i == 0 {
				us.shortcut.EXPECT().GetShort(firstShortURL, LenShortURL).Times(1).Return(fmt.Sprint(i))
			} else {
				us.shortcut.EXPECT().GetShort(fmt.Sprint(i-1), LenShortURL).Times(1).Return(fmt.Sprint(i))
			}

			us.rep.EXPECT().AddURL(fmt.Sprint(i), longURL).Times(1).Return(sr.ErrorURLAlreadyExists)
		}

		us.shortcut.EXPECT().GetShort(fmt.Sprint(8), LenShortURL).Times(1).Return(lastShortUrl)
		us.rep.EXPECT().AddURL(lastShortUrl, longURL).Times(1).Return(nil)

		t.NewStep("Тестирование")
		res, err := us.usc.GetShortURL(longURL)

		t.NewStep("Проверка результатов")
		t.Require().NoError(err)
		t.Require().Equal(lastShortUrl, res)
	})

	t.Run("Добавление длинного URL c ошибкой проверки существования созданного короткого", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		us.shortcut.EXPECT().GetShort(longURL, LenShortURL).Times(1).Return(shortURL)
		us.rep.EXPECT().GetLongURL(shortURL).Times(1).Return("", testError)

		t.NewStep("Тестирование")
		_, err := us.usc.GetShortURL(longURL)

		t.NewStep("Проверка результатов")
		t.Require().ErrorIs(err, testError)
	})

	t.Run("Добавление длинного URL c ошибкой проверки существования созданного короткого", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		us.shortcut.EXPECT().GetShort(longURL, LenShortURL).Times(1).Return(shortURL)
		us.rep.EXPECT().GetLongURL(shortURL).Times(1).Return("", testError)

		t.NewStep("Тестирование")
		_, err := us.usc.GetShortURL(longURL)

		t.NewStep("Проверка результатов")
		t.Require().ErrorIs(err, testError)
	})
}

func TestRunUsecaseTest(t *testing.T) {
	suite.RunSuite(t, new(UsecaseSuite))
}
