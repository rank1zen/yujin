package interal



type Error struct {
	orig error
	msg string
	code ErrorCode
}

type ErrorCode uint


func WrapErrorf(orig error, code ErrorCode, format string, a ...interface{}) error {
	return &Error{
		orig: orig,
		code: code,
		msg: fmt.Sprintf(format, a...),
	}
}

