package shares

type ExpiredShareError struct {
}

func (e *ExpiredShareError) Error() string {
	return "share is expired"
}
