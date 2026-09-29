package oi

import (
	"fmt"
	"io"
	"strings"
)

type InvalidOffset interface {
	error
	InvalidOffset() (offset int64, whence int)
}

func errInvalidOffset(offset int64, whence int, msg string) error {
	var e InvalidOffset = &internalInvalidOffset{
		offset:offset,
		whence:whence,
		msg:msg,
	}

	return e
}

func errInvalidOffsetf(offset int64, whence int, format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)

	return errInvalidOffset(offset, whence, msg)
}

type internalInvalidOffset struct {
	offset int64
	whence int
	msg    string
}

func (receiver internalInvalidOffset) Error() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "oi: Invalid Offset: offset=%d whence=%d", receiver.offset, receiver.whence)
	switch receiver.whence {
	case io.SeekStart:
		builder.WriteString(" (Seek Start)")
	case io.SeekCurrent:
		builder.WriteString(" (Seek Current)")
	case io.SeekEnd:
		builder.WriteString(" (Seek End)")
	}

	if "" != receiver.msg {
		builder.WriteRune(' ')
		builder.WriteString(receiver.msg)
	}

	return builder.String()
}

func (receiver internalInvalidOffset) InvalidOffset() (offset int64, whence int) {
	return receiver.offset, receiver.whence
}
