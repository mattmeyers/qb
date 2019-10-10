package qb

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrMissingTable = Error("no table specified")
)
