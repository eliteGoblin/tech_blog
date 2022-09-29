// built-in interface type

type error interface {
        Error() string
}

package errors
// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}
// errorString is a trivial implementation of error.
type errorString struct {
	s string
}
func (e *errorString) Error() string {
	return e.s
}
package fmt
func Errorf(format string, a ...interface{}) error {
	return errors.New(Sprintf(format, a...))
}