package shortener

// RedirectSerializer iterface is an adapter which
// connects the Redirect Service and HTTP request/response.
type RedirectSerializer interface {
	Decode(input []byte) (*Redirect, error)
	Encode(input *Redirect) ([]byte, error)
}
