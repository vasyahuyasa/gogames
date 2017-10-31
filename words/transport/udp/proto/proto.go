package proto

// CommandType строковое представление типа команды
type CommandType string

const (
	// CmdRegister клиент прислал команду о регистрации в игре
	CmdRegister CommandType = "register"

	// CmdSupply клиент прислал команду с запрашиваемым словом
	CmdSupply CommandType = "supply"

	// CmdDemand сервер рассылает команду с запрашиваемым словом
	CmdDemand CommandType = "demand"
)

// Cmd тип команды, используется для идентефикации
type Cmd struct {
	Command CommandType `json:"command"`
}

// Register заявка на рекистрацию от клиента
type Register struct {
	Cmd
	Name string `json:"name"`
}

// Deamand пакет отправляемый всем клиентам для запроса слова
type Deamand struct {
	Cmd            // demand
	Word    string `json:"word"`
	Letter  string `json:"letter"`
	Timeout int    `json:"timeout"`
	Player  string `json:"player"`
}

// Supply пакет отправляемый серверу с запрашиваемым словом
type Supply struct {
	Cmd         // supply
	Word string `json:"word"`
}

// Error ошибка передаваемая от сервера клиенту
type Error struct {
	Cmd         // error
	Code string `json:"code"`
	Msg  string `json:"msg"`
}
