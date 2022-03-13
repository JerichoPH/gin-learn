package errors

type UnAuthorizationError struct{ s string }

func (cls *UnAuthorizationError) Error() string {
	return cls.s
}

func ThrowUnAuthorization(text string) error {
	return &UnAuthorizationError{text}
}
