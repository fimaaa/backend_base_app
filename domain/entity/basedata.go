package entity

type BaseReqFind struct {
	Page   int
	Size   int
	Value  interface{}
	SortBy map[string]interface{}
}

type MyError struct {
	Message string
}

// Implement the Error method for MyError
func (e MyError) Error() string {
	return e.Message
}

// Define a function that returns a MyError
func NewMyError(message string) MyError {
	return MyError{Message: message}
}
