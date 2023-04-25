package domain

type IRedirectRepository interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
