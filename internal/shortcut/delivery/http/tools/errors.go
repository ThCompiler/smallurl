package tools

import (
	"smallurl/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Error string `json:"error,omitempty"`
}

// SendError завершает контекст запроса с переданной ошибкой и статус кодом.
//
// Также логгирует отправленный результат.
func SendError(c *gin.Context, err error, code int, l logger.Interface) {
	c.AbortWithStatusJSON(code, Error{Error: err.Error()})
	l.Info("[HTTP] Result - error %s was sent with status code %d", err, code)
}
