package apub

import (
	"errors"

	"golang.org/x/xerrors"
)

var (
	ErrLangNotFound     = errors.New("key not translated to given language")
	ErrLangMapNotFound  = errors.New("key has no language map")
	ErrKeyNotObject     = errors.New("unable to decode key as object")
	ErrKeyTypeNotObject = errors.New("unable to decode type as object")
	ErrInvalidInt       = errors.New("unable to decode value as int")
	ErrInvalidBool      = errors.New("unable to decode value as bool")
)

func FatalLangErr(err error) bool {
	if xerrors.Is(err, ErrLangNotFound) {
		return false
	}
	if xerrors.Is(err, ErrLangMapNotFound) {
		return false
	}
	return true
}
