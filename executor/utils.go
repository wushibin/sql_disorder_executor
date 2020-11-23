package executor

func BytesToString(b []byte) string {
	//return *((*string)(unsafe.Pointer(&b)))
	return string(b)
}
