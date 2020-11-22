package shortener

// RedirectService interface describes base
// functions of business logic.
type RedirectService interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
