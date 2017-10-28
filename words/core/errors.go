package core

// Неэкспортируемый тип ошибки. Удовлетворяет интерфейсу error
type e struct {
	text string
}

func (err e) Error() string {
	return err.text
}

// AlredyStarted возникает когда пытаются повторно запустить игру
var AlredyStarted = e{"Игра уже запущена"}

// AlreadyRegistered возникает когда игрок уже в списке
var AlreadyRegistered = e{"Игрок с таким именем уже зарегестрирован"}

// GameInProgress возникает когда игрок пытается присоедениться к запущеной игре
var GameInProgress = e{"Нельзя присоединяться к запущенной игре"}

// UnknownPlayer возникает когда игрок даёт ответ с неверной парой token и password
var UnknownPlayer = e{"Игрок не найден или не совпадает имя и пароль"}

// WrongPlayer возникает если сейчас другой игрок должен сделать ход
var WrongPlayer = e{"Ход другого игрока"}

// EmptyWord игрок прислал пустой ответ
var EmptyWord = e{"Пустой ответ"}

// MissedInDictonary игрок прислал неизвесное слово
var MissedInDictonary = e{"Слово не представлено в словаре"}

// UsedWord когда игрок присылает использованное слово
var UsedWord = e{"Слово уже было использоано"}

// Mistmatch возникает если игрок прислал слово которое начинается на другаю букву
var Mistmatch = e{"Присланное слово должно начинаться на другую букву"}
