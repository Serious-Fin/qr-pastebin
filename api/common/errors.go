package common

type NotFoundError struct {
}

func (e *NotFoundError) Error() string {
	return "resource not found"
}

type PasswordIncorrectError struct {
}

func (e *PasswordIncorrectError) Error() string {
	return "password is incorrect"
}
