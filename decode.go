package apub

import (
	"encoding/json"
	"io"
)

type Decoder struct {
	Language string
}

func (d *Decoder) Decode(input io.Reader) (*Object, error) {
	data := make(map[string]interface{})
	err := json.NewDecoder(input).Decode(&data)

	obj := New(data)
	if len(d.Language) > 0 {
		obj.lang = d.Language
	}

	return obj, err
}
