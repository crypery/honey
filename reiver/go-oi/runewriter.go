package oi

type RuneWriter interface {
	WriteRune(r rune) (n int, err error)
}
