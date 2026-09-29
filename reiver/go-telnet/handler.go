package telnet

import (
	"net"
)

type Handler interface {
	ServeTELNET(net.Conn,Context, Writer, Reader)
}
