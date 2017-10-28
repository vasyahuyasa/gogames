package udp

import (
	"io"
)

type client struct {
	Name string
	W    io.Writer
}

type clients []*client

// byName находит клиент а по имени
func (c clients) byName(name string) (client, bool) {
	for _, clt := range c {
		if clt.Name == name {
			return clt, true
		}
	}

	return nil, false
}
