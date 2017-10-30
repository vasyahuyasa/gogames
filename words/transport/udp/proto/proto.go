package proto

// CommandType строковое представление типа команды
type CommandType string

const (
	// CmdRegister клиент прислал команду о регистрации в игре
	CmdRegister CommandType = "register"

	// CmdSupply клиент прислал команду с запрашиваемым словом
	CmdSupply = "supply"

	// CmdDemand сервер рассылает команду с запрашиваемым словом
	CmdDemand = "demand"
)

// Command тип команды, используется для идентефикации
type Command struct {
	Command CommandType `json:"command"`
}

// Register заявка на рекистрацию от клиента
type Register struct {
	Command
	Name string `json:"name"`
}

// Deamand пакет отправляемый всем клиентам для запроса слова
type Deamand struct {
	Command        // demand
	Word    string `json:"word"`
	Letter  string `json:"letter"`
	Timeout int    `json:"timeout"`
	Player  string `json:"player"`
}

// Supply пакет отправляемый серверу с запрашиваемым словом
type Supply struct {
	Command        // supply
	Word    string `json:"word"`
}

// Error ошибка передаваемая от сервера клиенту
type Error struct {
	Command        // error
	Code    string `json:"code"`
	Msg     string `json:"msg"`
}
