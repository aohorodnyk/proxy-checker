package core

type proxyError struct {
	Message string
}

func (e proxyError) Error() string {
	return e.Message
}
