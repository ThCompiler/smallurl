package usecase

import (
	sr "smallurl/internal/shortcut/repository"

	"github.com/pkg/errors"
)

type GetShortcut func() Shortcut

const (
	LenShortURL = 10
)

type Usecase struct {
	rp       Repository
	shortcut GetShortcut
}

func NewUsecase(rp Repository, shortcut GetShortcut) *Usecase {
	return &Usecase{
		rp:       rp,
		shortcut: shortcut,
	}
}

func (u *Usecase) GetLongURL(shortURL string) (string, error) {
	return u.rp.GetLongURL(shortURL)
}

func (u *Usecase) GetShortURL(longURL string) (string, error) {
	shortcut := u.shortcut()

	shortURL := shortcut.GetShort(longURL, LenShortURL)

	foundURL, err := u.rp.GetLongURL(shortURL)

	switch {
	case err == nil && foundURL == longURL:
		return shortURL, nil
	case err == nil && foundURL != longURL:
		shortURL = shortcut.GetShort(shortURL, LenShortURL)
	case err != nil && !errors.Is(err, sr.ErrorURLNotFound):
		return "", errors.Wrapf(err, "with long url %s short url %s", foundURL, shortURL)
	}

	if err = u.rp.AddURL(shortURL, longURL); err == nil {
		return shortURL, nil
	}

	for errors.Is(err, sr.ErrorURLAlreadyExists) {
		shortURL = shortcut.GetShort(shortURL, LenShortURL)
		err = u.rp.AddURL(shortURL, longURL)
	}

	if err != nil {
		return "", errors.Wrapf(err, "with long url %s short url %s", longURL, shortURL)
	}

	return shortURL, nil
}
