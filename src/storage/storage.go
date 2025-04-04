package storage

import "errors"

var (
	ErrURLNotFound   = errors.New("url not found")
	ErrURLExists     = errors.New("url already exists")
	ErrALIASNotFound = errors.New("alias not found")
)
