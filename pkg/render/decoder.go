package render

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func Decode(r *http.Request, v interface{}) error {
	if err := DecodeJSON(r.Body, v); err != nil {
		return err
	}

	return nil
}

func DecodeJSON(r io.Reader, v interface{}) error {
	defer io.Copy(ioutil.Discard, r)
	return json.NewDecoder(r).Decode(v)
}
