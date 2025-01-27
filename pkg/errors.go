package pkg

import (
	"bytes"
	"strconv"
)

type ErrorCustom struct {
	Code        int
	Message     string
	Description string
}

// Error implements the error interface.
func (e *ErrorCustom) Error() string {
	var buffer bytes.Buffer
	buffer.Grow(60)
	buffer.WriteString("Message: ")
	buffer.WriteString(e.Message)
	buffer.WriteString(", Status Code: ")
	buffer.WriteString(strconv.Itoa(e.Code))
	buffer.WriteString(", Description : ")
	buffer.WriteString(e.Description)

	return buffer.String()
}
