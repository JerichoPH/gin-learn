package Errors

type ForbiddenError struct{ s string }

func (cls *ForbiddenError) Error() string {
	return cls.s
}

func ThrowForbidden(text string) error {
	return &ForbiddenError{text}
}
