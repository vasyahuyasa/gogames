package core

// Player объект игрока
type Player struct {
	// Имя для идентефикации, должно быть уникальным
	Name string

	// Флаг того, что игрок выбыл из игры
	out bool
}

// IsOut флаг того, выбыл ли игрок из игры
func (p *Player) IsOut() bool {
	return p.out
}

// Out устанавливает флаг того, что игрок выбыл из игры
func (p *Player) Out() {
	p.out = true
}
