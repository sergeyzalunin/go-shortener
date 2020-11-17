package shortener

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/ventu-io/go-shortid"
	"gopkg.in/validator.v2"
)

var (
	ErrRedirectNotFound = errors.New("Redirect Not Found")
	ErrRedirectInvalid  = errors.New("Redirect Invalid")
)

// redirectService an implementation of
// redirection business logic.
type redirectService struct {
	redirectRepo RedirectRepository
}

// NewRedirectService instantiates a new instance of RedirectService.
// The base implementation is redirectService struct.
func NewRedirectService(redirectRepo RedirectRepository) RedirectService {
	return &redirectService{
		redirectRepo,
	}
}

// Find returns the redirection reference or error.
// It independs on database due to RedirectRepository interface.
func (r *redirectService) Find(code string) (*Redirect, error) {
	return r.redirectRepo.Find(code)
}

// Store validates URL, assigns the short id to code and creation time.
// After that passes it to redirect repository to store in a db.
func (r *redirectService) Store(redirect *Redirect) error {
	// let's validate the correctness of url by provided rules in Redirect struct
	if err := validator.Validate(redirect); err != nil {
		// exchange the validator error for our stub and put the path it happens
		return errs.Wrap(ErrRedirectInvalid, "service.Redirect.Store")
	}

	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()

	return r.redirectRepo.Store(redirect)
}
