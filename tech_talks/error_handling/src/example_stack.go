func ExampleWrap_extended() {
	err := fn()
	fmt.Printf("%+v\n", err)
	// Example output:
	// error
    // _/home/frankie/git_repo/errors_test.fn
	//    /home/frankie/git_repo/errors/example_test.go:100
    // _/home/frankie/git_repo/errors_test.ExampleWrap_extended
	//    /home/frankie/git_repo/errors/example_test.go:116
	// ... more error stack
	// inner
	// ... inner's stack
	// middle
	// ... middle's stack
	// outter
	// ... outter's stack
}
