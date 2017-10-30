package core

// Player объект игрока
type Player struct {
	// Имя для идентефикации, должно быть уникальным
	Name string

	// Флаг того, что игрок выбыл из игры
	IsOut bool
}
