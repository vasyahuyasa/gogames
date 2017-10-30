package transport

import (
	"github.com/vasyahuyasa/gogames/words/core"
)

// Interface описывает общие команды для для взаимодействия с ядром
// через любой протокол
type Interface interface {
	// SendTurn отсылает всем игрокам следующее слово и игрока от которого ожидается ответ
	SendTurn(core.Turn) error

	// RegChan возвращает канал из которого будут приходить запросы о регистрации в матче
	RegChan() <-chan RegInfo

	// TurnChan возвращает канал из которого будут приходить ответы на запрошенные слова
	TurnChan() <-chan Turn

	// Отослать ошибку игроку если to пустое, то разослать всем
	Error(to string, err error)
}

// RegInfo информация посылаемая игроком для регистрации в перед началом
type RegInfo struct {
	Name string
}

// Turn информация песылаемая игроком для ответа на команду о запросе слова
type Turn struct {
	Name string
}
