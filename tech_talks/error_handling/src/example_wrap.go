func ExampleWrap() {
	cause := errors.New("file not exist")
	err := errors.Wrap(cause, "can not write content")
	fmt.Println(err)

	// Output: can not write content: file not exist
}