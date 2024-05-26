package repository

import "github.com/pkg/errors"

var (
	ErrorURLAlreadyExists = errors.New("short URL already exists")
	ErrorURLNotFound      = errors.New("short URL not found")
)
