package handlers

import "github.com/pkg/errors"

var (
	ErrorURLNotValid          = errors.New("original URL not valid")
	ErrorShortURLNotPresented = errors.New("short URL not present")

	ErrorServerError = errors.New("got system error, try again later")
)
