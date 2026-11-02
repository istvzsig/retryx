package retry

import "fmt"

type HTTPStatusError struct {
	StatusCode int
}

var _ error = HTTPStatusError{}

func (e HTTPStatusError) Error() string {
	return fmt.Sprintf("retry: http status %d", e.StatusCode)
}
