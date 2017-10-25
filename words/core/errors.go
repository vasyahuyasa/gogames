package core

// Неэкспортируемый тип ошибки. Удовлетворяет интерфейсу error
type e struct {
	text string
}

func (err e) Error() string {
	return err.text
}

// AlreadyRegistered возникает когда игрок уже в списке
var AlreadyRegistered = e{"Игрок с таким именем уже зарегестрирован"}

// GameInProgress возникает когда игрок пытается присоедениться к запущеной игре
var GameInProgress = e{"Игра уже началась"}

// AlredyStarted возникает когда пытаются повторно запустить игру
var AlredyStarted = e{"Игра уже запущена"}

// UnknownPlayer возникает когда игрок даёт ответ с неверной парой token и password
var UnknownPlayer = e{"Игрок не найден или не совпадает токен и пароль"}
