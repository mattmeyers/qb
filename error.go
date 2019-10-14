package qb

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrMissingTable          = Error("no table specified")
	ErrMissingSetPairs       = Error("no set pairs provided")
	ErrColValMismatch        = Error("the number of columns and values do not match")
	ErrInvalidConflictTarget = Error("invalid conflict target")
	ErrInvalidConflictAction = Error("invalid conflict action")
)
