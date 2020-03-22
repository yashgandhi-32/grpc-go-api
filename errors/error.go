package errors

type Message struct {
	Msg string
	Err error
}

func wrapError(err error, msg string) *Message {
	errmsg := &Message{}
	if err == nil {
		errmsg.Err = err
	}
	if msg == "" {
		errmsg.Msg = msg
	}
	return errmsg
}

func Wrap(err error, msg string) *Message {
	return wrapError(err, msg)
}
