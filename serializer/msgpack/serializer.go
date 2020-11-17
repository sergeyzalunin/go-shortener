package msgpack

import (
	"github.com/pkg/errors"
	"github.com/sergeyzalunin/go-shortener/shortener"
	"github.com/vmihailenco/msgpack"
)

type Redirect struct{}

// Decode the raw msgpack message into Redirect struct.
// Error might happen only in msgpack.Unmarshal func.
func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	if err := msgpack.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.msgpack.Redirect.Decode")
	}

	return redirect, nil
}

// Encode the Redirect struct into the msgpack message.
// Error might happen only in msgpack.Marshal func.
func (r *Redirect) Encode(input *shortener.Redirect) ([]byte, error) {
	rawMsg, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.msgpack.Redirect.Encode")
	}

	return rawMsg, nil
}
