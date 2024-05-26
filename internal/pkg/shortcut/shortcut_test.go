package shortcut

import (
	"crypto/md5"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
)

const (
	resultWordLen = 10
)

func TestShortcut(t *testing.T) {
	runner.Run(t, "Тестирование функции Shortcut", func(t provider.T) {
		t.NewStep("Инициализация исходных данных")
		inputString := "Какая-то случайная строка"

		t.WithNewStep("Успешное получение сокращённой версии входных данных", func(t provider.StepCtx) {
			t.NewStep("Инициализация Shortcut")
			shrt := NewHashShortcut()

			t.NewStep("Тестирование")
			res := shrt.GetShort(inputString, resultWordLen)

			t.NewStep("Проверка результата")
			t.Require().Len(res, resultWordLen)
			for _, c := range res {
				t.Require().Contains([]byte(alphabet), byte(c))
			}
		})

		t.WithNewStep("Ошибочное получение сокращённой версии длинной большей, чем позволяет используемая хэш функция",
			func(t provider.StepCtx) {
				t.NewStep("Инициализация Shortcut")
				shrt := NewHashShortcut()

				defer func() {
					if err := recover(); err != nil {
						e, ok := err.(error)
						t.Require().True(ok)
						t.Require().Equal(e, ErrorBadHashFunction)
						return
					}
				}()
				t.NewStep("Тестирование")
				shrt.GetShort(inputString, md5.Size+1)

				t.Require().False(true, "Функция не среагировала на ошибочный размер сокращённого слова")
			})
	})
}
