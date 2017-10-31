package udp

import (
	"io"
	"net"
)

type client struct {
	name string
	addr *net.UDPAddr
	W    io.Writer
}

type clients []*client

// byName находит клиент а по имени
func (c *clients) byName(name string) (*client, bool) {
	for _, clt := range *c {
		if clt.name == name {
			return clt, true
		}
	}

	return nil, false
}

// byAddr находит клиента по адресу
func (c *clients) byAddr(addr *net.UDPAddr) (*client, bool) {
	for _, clt := range *c {
		if clt.addr.String() == addr.String() {
			return clt, true
		}
	}

	return nil, false
}
