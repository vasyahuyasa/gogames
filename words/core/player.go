package core

// Player объект игрока
type Player struct {
	// Имя для идентефикации, должно быть уникальным
	Name string

	// Пароль для верификации ответа
	Password string

	// Токен по которому можно идентефецировать игрока
	Token string

	// Флаг того, что игрок выбыл из игры
	IsOut bool
}
