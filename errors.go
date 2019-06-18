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
	ErrInvalidBool      = errors.New("unable to decode value as bool")
	ErrInvalidFloat     = errors.New("unable to decode value as float")
	ErrInvalidInt       = errors.New("unable to decode value as int")
	ErrInvalidTime      = errors.New("unable to decode value as time")
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
