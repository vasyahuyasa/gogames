package core

import (
	"strings"
	"time"
	"unicode/utf8"
)

// Длина случайно сгенерированого токена
const tokenLen = 12

// Буквы на которые не может начинаться следующее слово
const excludeLetters = "цфзшэЪЬЫЙ"

// Game отвечает за всю игровую логику
type Game struct {
	players []*Player

	started     bool
	turnTimeout time.Duration
	turnChan    chan Turn

	word       string
	nextRune   string
	nextPlayer int
}

// findPlayer если находит игрока возвращает объект, индекс и флаг ok = true
func (g *Game) findPlayer(token string, password string) (*Player, int, bool) {
	for i, p := range g.players {
		if p.Token == token && p.Password == password {
			return p, i, true
		}
	}
	return nil, 0, false
}

// isRegistered проверять есть ли игрок с именем name в списке
func (g *Game) isRegistered(name string) bool {
	for _, p := range g.players {
		if p.Name == name {
			return true
		}
	}
	return false
}

// validWord проверяет является ли присланное игроком слово верным
func (g *Game) validWord(word string) bool {
	if len(g.nextRune) == 0 {
		return true
	}

	// разделить присланное слово на отдельные буквы
	wordRunes := strings.Split(strings.ToLower(word), "")

	return wordRunes[0] == g.nextRune
}

// nextWord устанавливает текущее слово и слудующую букву
func (g *Game) nextWord(word string) {
	g.word = strings.Title(word)
	g.nextRune = ""

	// разделить присланное слово на отдельные буквы
	wordRunes := strings.Split(strings.ToLower(word), "")

	// пройти по массиву букв начиная с последней и найти ту,
	// с которой будет начинаться следующее слово
	for i := len(wordRunes) - 1; i >= 0; i-- {
		if !strings.Contains(excludeLetters, wordRunes[i]) {
			g.nextRune = wordRunes[i]
			break
		}
	}

}

// findNextPlayer ищет index следующего не выбывшего игрока, если не найден возвращает ok = false
func (g *Game) findNextPlayer() (int, bool) {
	currentindex := g.nextPlayer

	// оставшиеся игроки до конца массива
	for i, p := range g.players[g.nextPlayer+1:] {
		if !p.IsOut {
			return i, true
		}
	}

	// игроки с начала массива до текущего
	for i, p := range g.players[:g.nextPlayer] {
		if !p.IsOut {
			return i, true
		}
	}

	return 0, false
}

// RegisterPlayer добавляет игрока, игрок может быть добавлен только до начала партии.
// Игрок передаёт свой пароль и взамен получает уникальный токен.
func (g *Game) RegisterPlayer(name string, password string) (string, error) {
	if g.isRegistered(name) {
		return "", AlreadyRegistered
	}

	token := uniqueID(tokenLen)
	p := &Player{
		Name:     name,
		Password: password,
		Token:    token,
	}
	g.players = append(g.players, p)

	return token, nil
}

// MakeTurn принимает ответ от игрока, если такого игрока нет в списке возвращается ошибка UnknownPlayer.
// Пустой word считается за проигрыш.
// Не совпадение последней буквы и первой буквы присланного слова считается за проигрыш.
func (g *Game) MakeTurn(token string, password string, word string) error {
	p, i, ok := g.findPlayer(token, password)
	if !ok {
		return UnknownPlayer
	}

	// пустой ответ
	if utf8.RuneCountInString(word) == 0 {
		p.IsOut = true
		g.nextTurn()
		return nil
	}

}

// Start начинает игру
func (g *Game) Start() (<-chan Turn, error) {
	if g.started {
		return nil, AlredyStarted
	}

	g.turnChan = make(chan Turn, 1000)
	g.started = true
	return g.turnChan, nil
}
