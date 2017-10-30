package udp

import (
	"encoding/json"
	"fmt"
	"github.com/vasyahuyasa/gogames/words/core"
	"github.com/vasyahuyasa/gogames/words/transport"
	"github.com/vasyahuyasa/gogames/words/transport/udp/proto"
	"log"
	"net"
)

const bufSize = 1024 * 8 // 8 KB
const chanSize = 1000

// Service иплементирует transport.Interface поверх протокола UDP
type Service struct {
	clients clients

	regs  chan transport.RegInfo
	turns chan transport.Turn
}

// SendTurn отправляет пакет с запросом слова всем подключенным клиентам индивидуально
func (s *Service) SendTurn(core.Turn) error {
	return nil
}

// RegChan врзвращает канал через который будут поступать заявки на регистрацию
func (s *Service) RegChan() <-chan transport.RegInfo {
	return s.regs
}

// Error отправляет ошибку клиенту, если клиент не указан, то всем
func (s *Service) Error(to string, err error) {

}

// parseData обрабатывае входные данные от клиента
func (s *Service) parseData(from *net.UDPAddr, data []byte) error {
	c := proto.Command{}
	err := json.Unmarshal(data, &c)
	if err != nil {
		return err
	}

	switch c.Command {

	// Клиент регистрируется в игре
	case proto.CmdRegister:
		r := proto.Register{}
		err = json.Unmarshal(data, &r)
		if err != nil {
			return err
		}

		s.regs <- transport.RegInfo{
			Name: r.Name,
		}

	// Клиент отправил слово
	case proto.CmdSupply:
		sup := proto.Supply{}
		err = json.Unmarshal(data, &sup)
		if err != nil {
			return err
		}

		name, ok := s.clients.byAddr(from)
		if !ok {
			return fmt.Errorf("Клиент не зарегистрирован: %s", from)
		}

		s.turns <- transport.Turn{
			Name: name,
		}

	}

	return fmt.Errorf("Неизвестная команда: %q", c.Command)
}

// run запускает фоновый udp listener
func (s *Service) run(conn *net.UDPConn) {
	go func() {
		buf := make([]byte, bufSize)

		for {
			n, addr, err := conn.ReadFromUDP(buf)

			if err != nil {
				log.Printf("ReadFromUDP: %v", err)
				continue
			}

			s.parseData(addr, buf[:n])
		}
	}()
}

// New создаёт новый экземпляр UDP транспорта
// параметром передаётся порт на котором будет висеть сервис
func New(port uint) (*Service, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	s := &Service{
		regs:  make(chan transport.RegInfo, chanSize),
		turns: make(chan transport.Turn, chanSize),
	}
	s.run(conn)

	return s, nil
}
