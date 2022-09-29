type MyError struct {
	error
}

func (err MyError) Error() string {
	return err.Error()
}