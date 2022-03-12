package Errors

type EmptyError struct{ s string }

func (cls *EmptyError) Error() string {
	return cls.s
}

func ThrowEmpty(text string) error {
	return &EmptyError{text}
}
