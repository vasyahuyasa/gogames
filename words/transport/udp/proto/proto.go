package proto

// Command тип приходящей или исходящей команды
type Command struct {
	Command string `json:"command"`
}

// Register заявка на рекистрацию от клиента
type Register struct {
	Command
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Turn пакет отправляемый всем клиентам для запроса слова
type Turn struct {
	Command        // demand
	Word    string `json:"word"`
	Letter  string `json:"letter"`
	Timeout int    `json:"timeout"`
	Player  string `json:"player"`
}
