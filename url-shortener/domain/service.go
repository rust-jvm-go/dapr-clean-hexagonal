package domain

type IRedirectService interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
