package shares

type PasswordIncorrectError struct {
}

func (e *PasswordIncorrectError) Error() string {
	return "password is incorrect"
}
