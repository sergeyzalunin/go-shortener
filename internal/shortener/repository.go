package shortener

// RedirectRepository interface describes
// base functions to database requests.
type RedirectRepository interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
