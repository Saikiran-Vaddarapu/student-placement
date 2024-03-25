package customerrors

type InternalErrors struct {
	Message string
}

func (e InternalErrors) Error() string {
	return e.Message
}
