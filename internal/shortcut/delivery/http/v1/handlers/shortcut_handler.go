package handlers

import (
	"net/http"
	"net/url"
	"smallurl/internal/app/delivery/http/middleware"
	"smallurl/internal/shortcut/delivery"
	"smallurl/internal/shortcut/delivery/http/tools"
	"smallurl/internal/shortcut/delivery/http/v1/models/request"
	"smallurl/internal/shortcut/delivery/http/v1/models/response"
	sr "smallurl/internal/shortcut/repository"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	ShortURLParam = "short_url"
)

type ShortcutHandlers struct {
	usc delivery.Usecase
}

func NewShortcutHandlers(usc delivery.Usecase) *ShortcutHandlers {
	return &ShortcutHandlers{usc: usc}
}

// GenShort
//
//	@Summary		Сокращение оригинального URL.
//	@Description	Сохраняет оригинальный URL в базе и возвращает сокращённый.
//	@Tags			urls
//	@Accept			json
//	@Param			request	body	request.LongURL	true	"Значение оригинального URL"
//	@Produce		json
//	@Success		201	{object}	response.Result	"Оригинальный URL успешно добавлен в систему"
//	@Failure		400	{object}	tools.Error		"Некорректные данные запроса"
//	@Failure		500	{object}	tools.Error		"Внутренняя ошибка сервера"
//	@Router			/shorten [post]
func (sh *ShortcutHandlers) GenShort(c *gin.Context) {
	l := middleware.GetLogger(c)

	var longURL request.LongURL
	if code, err := tools.ParseRequestBody(c.Request.Body, &longURL, request.ValidateLongURL, l); err != nil {
		tools.SendError(c, err, code, l)

		return
	}

	// Проверка валидности URL
	if _, err := url.ParseRequestURI(longURL.OriginalURL); err != nil {
		tools.SendError(c, ErrorURLNotValid, http.StatusBadRequest, l)

		return
	}

	shortURL, err := sh.usc.GetShortURL(longURL.OriginalURL)
	if err != nil {
		tools.SendError(c, ErrorServerError, http.StatusInternalServerError, l)
		l.Error("[HTTP} - %s", errors.Wrapf(err, "can't create short url"))

		return
	}

	tools.SendStatus(c, http.StatusCreated, &response.Result{ShortURL: shortURL}, l)
}

// GetLong
//
//	@Summary		Получение оригинального URL.
//	@Description	Принимает сокращённый URL и возвращает оригинальный.
//	@Tags			urls
//	@Param			short_url	path	string	true	"Короткий URL"
//	@Produce		json
//	@Success		303	"Найден оригинальный URL, и запрос перенаправлен на него"
//	@Failure		400	{object}	tools.Error	"Некорректные данные запроса"
//	@Failure		404	"Оригинальный URL не найден"
//	@Failure		500	{object}	tools.Error	"Внутренняя ошибка сервера"
//	@Router			/{short_url} [get]
func (sh *ShortcutHandlers) GetLong(c *gin.Context) {
	l := middleware.GetLogger(c)

	// Получение ключа
	shortURL := c.Param(ShortURLParam)
	if shortURL == "" {
		tools.SendError(c, ErrorShortURLNotPresented, http.StatusBadRequest, l)

		return
	}

	longURL, err := sh.usc.GetLongURL(shortURL)
	if err != nil {
		if errors.Is(err, sr.ErrorURLNotFound) {
			tools.SendStatus(c, http.StatusNotFound, nil, l)

			return
		}

		tools.SendError(c, ErrorServerError, http.StatusInternalServerError, l)

		l.Error("[HTTP} - %s", errors.Wrapf(err, "can't get original URL from short URL"))

		return
	}

	c.Redirect(http.StatusSeeOther, longURL)
}
