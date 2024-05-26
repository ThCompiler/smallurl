package memory

import (
	"smallurl/internal/shortcut/repository"
	"sync"

	"github.com/pkg/errors"
)

type Repository struct {
	mem map[string]string
	mc  sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{mem: make(map[string]string)}
}

func (rp *Repository) GetLongURL(shortURL string) (string, error) {
	rp.mc.RLock()
	defer rp.mc.RUnlock()

	if longURL, ok := rp.mem[shortURL]; ok {
		return longURL, nil
	}

	return "", errors.Wrapf(repository.ErrorURLNotFound, "short url %s not found", shortURL)
}

func (rp *Repository) AddURL(shortURL, longURL string) error {
	rp.mc.Lock()
	defer rp.mc.Unlock()

	if _, ok := rp.mem[shortURL]; ok {
		return errors.Wrapf(repository.ErrorURLAlreadyExists,
			"try add long url %s with short url %s", longURL, shortURL)
	}

	rp.mem[shortURL] = longURL

	return nil
}
