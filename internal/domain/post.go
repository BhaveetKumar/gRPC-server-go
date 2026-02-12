package domain

import "github.com/BhaveetKumar/gRPC-server-go/internal/errors"

type Post struct {
	ID              string
	Title           string
	Content         string
	Author          string
	PublicationDate string
	Tags            []string
}

func (p *Post) Validate() error {
	if p == nil {
		return errors.ErrInvalidInput
	}

	if p.Title == "" || p.Content == "" || p.Author == "" {
		return errors.ErrInvalidInput
	}

	return nil
}
