package core

// coreError - общий тип ошибок не фатальна для игрока
type coreError struct {
	text string
}

// playerError - если игрок получает этот тип ошибки, то он выбывает из игры
type playerError struct {
	text string
}

func (err coreError) Error() string {
	return err.text
}

func (err playerError) Error() string {
	return err.text
}

// NotEnoughtPlayers возникает когда вызывается команда start с недостаточным количеством игроков
var NotEnoughtPlayers = coreError{"Недостаточно игроков чтобы начать"}

// AlredyStarted возникает когда пытаются повторно запустить игру
var AlredyStarted = coreError{"Игра уже запущена"}

// AlreadyRegistered возникает когда игрок уже в списке
var AlreadyRegistered = coreError{"Игрок с таким именем уже зарегестрирован"}

// GameInProgress возникает когда игрок пытается присоедениться к запущеной игре
var GameInProgress = coreError{"Нельзя присоединяться к запущенной игре"}

// UnknownPlayer возникает когда имя игрока не зарегестрировано
var UnknownPlayer = coreError{"Игрок не найден"}

// WrongPlayer возникает если сейчас другой игрок должен сделать ход
var WrongPlayer = coreError{"Ход другого игрока"}

// EmptyWord игрок прислал пустой ответ
var EmptyWord = playerError{"Пустой ответ"}

// MissedInDictonary игрок прислал неизвесное слово
var MissedInDictonary = playerError{"Слово не представлено в словаре"}

// UsedWord когда игрок присылает использованное слово
var UsedWord = playerError{"Слово уже было использоано"}

// Mistmatch возникает если игрок прислал слово которое начинается на другаю букву
var Mistmatch = playerError{"Присланное слово должно начинаться на другую букву"}
