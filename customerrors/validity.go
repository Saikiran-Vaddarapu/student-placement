package customerrors

type Validity struct {
	Message string
}

func (v Validity) Error() string {
	return v.Message
}
