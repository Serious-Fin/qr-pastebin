package users

type UserAlreadyExistsError struct {
}

func (e *UserAlreadyExistsError) Error() string {
	return "user already exists"
}

type WrongPasswordError struct {
}

func (e *WrongPasswordError) Error() string {
	return "username or password is incorrect"
}
