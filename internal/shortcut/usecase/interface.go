package usecase

//go:generate mockgen -destination=mocks/repository.go -package=mu . Repository

type Repository interface {
	GetLongURL(shortURL string) (string, error)
	AddURL(shortURL, longURL string) error
}

//go:generate mockgen -destination=mocks/shortcut.go -package=mu . Shortcut

type Shortcut interface {
	GetShort(s string, n int) string
}
