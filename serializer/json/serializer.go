package json

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/sergeyzalunin/go-shortener/shortener"
)

type Redirect struct{}

// Decode the raw json message into Redirect struct.
// Error might happen only in json.Unmarshal func.
func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Json.Redirect.Decode")
	}

	return redirect, nil
}

// Encode the Redirect struct into the json message.
// Error might happen only in json.Marshal func.
func (r *Redirect) Encode(input *shortener.Redirect) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Json.Redirect.Encode")
	}

	return rawMsg, nil
}
