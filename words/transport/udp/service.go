package udp

import (
	"github.com/vasyahuyasa/gogames/words/core"
)

// Service иплементирует transport.Interface повер протокола UDP
type Service struct {
	clients clients
}

// SendTurn отправляет пакет с запросом слова всем подключенным клиентам индивидуально
func (s *Service) SendTurn(core.Turn) error {

}

// RegChan врзвращает канал через который будут поступать заявки на регистрацию
func (s *Service) RegChan() <-chan RegInfo {

}

// Error отправляет ошибку клиенту, если клиент не указан, то всем
func (s *Service) Error(to string, err error) {

}
