package shares

import "fmt"

type PasswordIncorrectError struct {
}

func (e *PasswordIncorrectError) Error() string {
	return fmt.Sprintf("password is incorrect")
}
