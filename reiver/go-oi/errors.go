package oi

import (
	"errors"
)

var (
	errNilReaderAt = errors.New("oi. Nil io.ReaderAt")
	errNilReceiver = errors.New("oi: Nil Receiver")
)
