package test

import (
	"bytes"
	"encoding/json"
	"io"
)

func GetBody(body map[string]interface{}) io.Reader {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(body)
	return b
}
