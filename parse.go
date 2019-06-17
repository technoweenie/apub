package apub

import (
	"encoding/json"
	"io"
)

type Parser struct {
	Language string
}

func (p *Parser) Parse(input io.Reader) (*Object, error) {
	data := make(map[string]interface{})
	err := json.NewDecoder(input).Decode(&data)

	obj := New(data)
	if len(p.Language) > 0 {
		obj.lang = p.Language
	}

	return obj, err
}
