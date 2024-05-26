package delivery

//go:generate mockgen -destination=mocks/usecase.go -package=md . Usecase

type Usecase interface {
	GetLongURL(shortURL string) (string, error)
	GetShortURL(longURL string) (string, error)
}
