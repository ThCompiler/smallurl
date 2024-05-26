package memory

import (
	"fmt"
	"smallurl/internal/shortcut/repository"
	"sync"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type MemoryRepositorySuite struct {
	suite.Suite
	rep *Repository
}

func (mrs *MemoryRepositorySuite) BeforeEach(t provider.T) {
	mrs.rep = NewRepository()
}

func (mrs *MemoryRepositorySuite) TestGetLongURL(t provider.T) {
	t.Title("Тестирование метода GetLongURL")

	t.Run("Успешное получение длинного URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"
		mrs.rep.mc.Lock()
		mrs.rep.mem[shortURL] = longURL
		mrs.rep.mc.Unlock()

		t.NewStep("Тестирование")
		res, err := mrs.rep.GetLongURL(shortURL)

		t.NewStep("Проверка результатов")
		t.Require().NoError(err)
		t.Require().Equal(longURL, res)
	})

	t.Run("Получение длинного URL по не существующему короткому", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		shortURL := "asdsvasdase"

		t.NewStep("Тестирование")
		_, err := mrs.rep.GetLongURL(shortURL)

		t.NewStep("Проверка результатов")
		t.Require().ErrorIs(err, repository.ErrorURLNotFound)
	})
}

func (mrs *MemoryRepositorySuite) TestAddURL(t provider.T) {
	t.Title("Тестирование метода AddURL")

	t.Run("Успешное добавление короткого URL с длинным URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"

		t.NewStep("Тестирование")
		err := mrs.rep.AddURL(shortURL, longURL)

		t.NewStep("Проверка результатов")
		t.Require().NoError(err)
	})

	t.Run("Добавление уже существующего короткого URL с длинным URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdase"

		mrs.rep.mc.Lock()
		mrs.rep.mem[shortURL] = longURL
		mrs.rep.mc.Unlock()

		t.NewStep("Тестирование")
		err := mrs.rep.AddURL(shortURL, longURL)

		t.NewStep("Проверка результатов")
		t.Require().ErrorIs(err, repository.ErrorURLAlreadyExists)
	})
}

func (mrs *MemoryRepositorySuite) TestRace(t provider.T) {
	t.Title("Тестирование методов AddURL и GetLongURL при записи и чтение из нескольких горутин")
	t.NewStep("Инициализация тестовых данных")
	longURL := "http://example.com"

	t.NewStep("Запуск горутин")
	wait := sync.WaitGroup{}
	wait.Add(4)

	go func() {
		for i := 0; i < 100; i++ {
			_ = mrs.rep.AddURL(fmt.Sprint(i), longURL)
		}
		wait.Done()
	}()

	go func() {
		for i := 0; i < 100; i++ {
			_ = mrs.rep.AddURL(fmt.Sprint(i+100), longURL)
		}
		wait.Done()
	}()

	go func() {
		for i := 0; i < 100; i++ {
			_, _ = mrs.rep.GetLongURL(fmt.Sprint(i))
		}
		wait.Done()
	}()

	go func() {
		for i := 0; i < 100; i++ {
			_, _ = mrs.rep.GetLongURL(fmt.Sprint(i + 100))
		}
		wait.Done()
	}()

	t.NewStep("Ожидание горутин")
	wait.Wait()

	t.NewStep("Проверка хранилища")
	for i := 0; i < 200; i++ {
		mrs.rep.mc.RLock()
		t.Require().Contains(mrs.rep.mem, fmt.Sprint(i))
		mrs.rep.mc.RUnlock()
	}
}

func TestRunMemoryRepositoryTest(t *testing.T) {
	suite.RunSuite(t, new(MemoryRepositorySuite))
}
