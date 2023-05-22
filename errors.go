package main

import "fmt"

const (
	GenericError = iota + 100 // generic DBS error
	BadRequest                // 101 bad request
	JsonMarshal               // 102 json.Marshal error
	FileIOError               // 103 file IO error
	SessionError              // 104 session error
)

// helper function to return human error message for given error code
func errorMessage(code int) string {
	if code == 0 {
		return ""
	} else if code == 101 {
		return "bad request"
	} else if code == 102 {
		return "JSON marshal error"
	} else if code == 103 {
		return "file IO error"
	} else if code == 104 {
		return "Session error"
	} else {
		return fmt.Sprintf("Not Implemented error for code %d", code)
	}
}
