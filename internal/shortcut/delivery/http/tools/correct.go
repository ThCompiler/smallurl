package tools

import (
	"smallurl/pkg/logger"

	"github.com/gin-gonic/gin"
)

// SendStatus завершает контекст запроса с переданным статус кодом и телом запроса. Если тело задано значением nil,
// будет отправлен только статус код.
//
// Также метод логирует отправленный статус код.
func SendStatus(c *gin.Context, code int, data any, l logger.Interface) {
	l.Info("[HTTP] Result - was sent response with status code %d", code)

	if data != nil {
		c.AbortWithStatusJSON(code, data)
	}

	c.AbortWithStatus(code)
}
