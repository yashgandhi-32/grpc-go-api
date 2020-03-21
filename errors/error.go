package errors

type message struct {
	msg string
	err error
}

func wrapError(err error, msg string) *message {
	errmsg := &message{}
	if err == nil {
		errmsg.err = err
	}
	if msg == "" {
		errmsg.msg = msg
	}
	return errmsg
}

func Wrap(err error, msg string) *message {
	return wrapError(err, msg)
}
