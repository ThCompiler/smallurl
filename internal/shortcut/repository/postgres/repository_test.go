package postgres

import (
	"smallurl/internal/shortcut/repository"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/pkg/errors"
)

var testError = errors.New("test error")

type PostgresRepositorySuite struct {
	suite.Suite
	mockPool pgxmock.PgxPoolIface
	rep      *Repository
}

func (prs *PostgresRepositorySuite) BeforeEach(t provider.T) {
	var err error
	prs.mockPool, err = pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}

	prs.rep = NewRepository(prs.mockPool)
}

func (prs *PostgresRepositorySuite) AfterEach(t provider.T) {
	prs.mockPool.Close()
}

func (prs *PostgresRepositorySuite) TestGetLongURL(t provider.T) {
	t.Title("Тестирование метода GetLongURL")

	t.Run("Успешное получение длинного URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		prs.mockPool.ExpectQuery(getLongURLQuery).WithArgs(shortURL).
			WillReturnRows(pgxmock.NewRows([]string{"long_url"}).AddRow(longURL)).Times(1)

		t.NewStep("Тестирование")
		res, err := prs.rep.GetLongURL(shortURL)

		t.NewStep("Проверка результатов")
		t.Require().NoError(err)
		t.Require().Equal(longURL, res)
	})

	t.Run("Получение длинного URL по не существующему короткому", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		prs.mockPool.ExpectQuery(getLongURLQuery).WithArgs(shortURL).
			WillReturnRows(pgxmock.NewRows([]string{"long_url"})).Times(1)

		t.NewStep("Тестирование")
		_, err := prs.rep.GetLongURL(shortURL)

		t.NewStep("Проверка результатов")
		t.Require().ErrorIs(err, repository.ErrorURLNotFound)
	})

	t.Run("Получение длинного URL с ошибкой", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		prs.mockPool.ExpectQuery(getLongURLQuery).WithArgs(shortURL).Times(1).WillReturnError(testError)

		t.NewStep("Тестирование")
		_, err := prs.rep.GetLongURL(shortURL)

		t.NewStep("Проверка результатов")
		t.Require().ErrorIs(err, testError)
	})
}

func (prs *PostgresRepositorySuite) TestAddURL(t provider.T) {
	t.Title("Тестирование метода AddURL")

	t.Run("Успешное добавление короткого URL с длинным URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		prs.mockPool.ExpectExec(addURLQuery).WithArgs(shortURL, longURL).
			WillReturnResult(pgxmock.NewResult("", 1)).Times(1)

		t.NewStep("Тестирование")
		err := prs.rep.AddURL(shortURL, longURL)

		t.NewStep("Проверка результатов")
		t.Require().NoError(err)
	})

	t.Run("Добавление уже существующего короткого URL с длинным URL", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		prs.mockPool.ExpectExec(addURLQuery).WithArgs(shortURL, longURL).
			Times(1).WillReturnError(&pgconn.PgError{
			Code:           uniqueConflictCode,
			ConstraintName: uniqueConstraintName,
		})

		t.NewStep("Тестирование")
		err := prs.rep.AddURL(shortURL, longURL)

		t.NewStep("Проверка результатов")
		t.Require().ErrorIs(err, repository.ErrorURLAlreadyExists)
	})

	t.Run("Добавление короткого URL с длинным URL с ошибкой", func(t provider.T) {
		t.NewStep("Инициализация тестовых данных")
		longURL := "http://example.com"
		shortURL := "asdsvasdasd"

		t.NewStep("Инициализация мока")
		prs.mockPool.ExpectExec(addURLQuery).WithArgs(shortURL, longURL).Times(1).WillReturnError(testError)

		t.NewStep("Тестирование")
		err := prs.rep.AddURL(shortURL, longURL)

		t.NewStep("Проверка результатов")
		t.Require().ErrorIs(err, testError)
	})
}

func TestRunPostgresRepositoryTest(t *testing.T) {
	suite.RunSuite(t, new(PostgresRepositorySuite))
}
