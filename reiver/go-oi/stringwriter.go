package oi

type StringWriter interface {
	WriteString(s string) (n int, err error)
}
