package errors

type UnLoginError struct{ s string }

func (cls *UnLoginError) Error() string {
	return cls.s
}

func ThrowUnLogin(text string) error {
	return &UnLoginError{text}
}
