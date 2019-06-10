package apubencoding

import (
	"fmt"
	"strconv"
)

const DefaultLang = "en"

type Object struct {
	lang string
	data map[string]interface{}
}

func (o *Object) String(name string) string {
	ival, ok := o.data[name]
	if !ok {
		return ""
	}

	switch val := ival.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func (o *Object) Context() string {
	return o.String("@context")
}

func (o *Object) ID() string {
	return o.String("id")
}

func (o *Object) Type() string {
	return o.String("type")
}

func (o *Object) Name() string {
	return o.String("name")
}
