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
	conn    *net.UDPConn
	regs    chan transport.RegInfo
	turns   chan transport.Turn
}

// SendTurn отправляет пакет с запросом слова всем подключенным клиентам индивидуально
func (s *Service) SendTurn(t core.Turn) error {
	cmd := proto.Deamand{
		Cmd: proto.Cmd{
			Command: proto.CmdDemand,
		},
		Word:    t.Word,
		Letter:  t.Letter,
		Timeout: int(t.Timeout.Seconds()),
		Player:  t.Player,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	for _, client := range s.clients {
		s.conn.WriteToUDP(data, client.addr)
	}

	return nil
}

// RegChan врзвращает канал через который будут поступать заявки на регистрацию
func (s *Service) RegChan() <-chan transport.RegInfo {
	return s.regs
}

// TurnChan возвращает канал через который клиенты присылают свои ответы
func (s *Service) TurnChan() <-chan transport.Turn {
	return s.turns
}

// Error отправляет ошибку клиенту, если клиент не указан, то всем
func (s *Service) Error(to string, err error) {

}

// parseData обрабатывае входные данные от клиента
func (s *Service) parseData(from *net.UDPAddr, data []byte) error {
	c := proto.Cmd{}
	err := json.Unmarshal(data, &c)
	if err != nil {
		return err
	}

	log.Println("Пакет:", c)

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

		c, ok := s.clients.byAddr(from)
		if !ok {
			return fmt.Errorf("Клиент не зарегистрирован: %s", from)
		}

		s.turns <- transport.Turn{
			Name: c.name,
			Word: sup.Word,
		}

	default:
		return fmt.Errorf("Неизвестная команда: %q", c.Command)
	}

	return nil
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

			err = s.parseData(addr, buf[:n])
			if err != nil {
				log.Printf("Ошибка parseData: %v", err)
			}
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
	log.Printf("Сервер запущен на %s", addr)

	s := &Service{
		conn:  conn,
		regs:  make(chan transport.RegInfo, chanSize),
		turns: make(chan transport.Turn, chanSize),
	}
	s.run(conn)

	return s, nil
}
