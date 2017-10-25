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
