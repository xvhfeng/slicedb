package eos

func IsEOF(err error) bool {
	switch pe := err.(type) {
	case nil:
		return false
	case *os.EOF:
		return true
	}
	return false
}
