package errors

import (
	pkgerrors "github.com/pkg/errors"
)

func New(message string) error {
	return pkgerrors.New(message)
}

func Annotate(other error, msg string) error {
	return pkgerrors.Wrap(other, msg)
}
