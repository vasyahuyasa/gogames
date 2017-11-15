// Клиент для тестирования сервера wordsd

package main

import (
	"net"
)

const address = "localhost:48879"

func main() {
	ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
}
