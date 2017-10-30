package core

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/vasyahuyasa/gogames/words/core/dictonary"
)

// Буквы на которые не может начинаться следующее слово
const excludeLetters = "цфзшэЪЬЫЙ"

// Game отвечает за всю игровую логику
type Game struct {
	players   []*Player
	dictonary *dictonary.Сollection
	used      *dictonary.Сollection

	started     bool
	turnTimeout time.Duration
	turnChan    chan Turn

	word       string
	nextRune   string
	nextPlayer int
}

// findPlayer если находит игрока возвращает объект, индекс и флаг ok = true
func (g *Game) findPlayer(name string) (*Player, int, bool) {
	for i, p := range g.players {
		if p.Name == name {
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
	if g.nextRune == "" {
		return true
	}

	// разделить присланное слово на отдельные буквы
	wordRunes := strings.Split(strings.ToLower(word), "")

	return wordRunes[0] == g.nextRune
}

// lastRune возвращает букву на которую должно начинаться следующеес слово
// функция отбрасывает запрещенные буквы
func (g *Game) lastRune(word string) string {
	last := ""

	// разделить присланное слово на отдельные буквы
	wordRunes := strings.Split(strings.ToLower(word), "")

	// пройти по массиву букв начиная с последней и найти ту,
	// с которой будет начинаться следующее слово игнорируя запрещенные символы
	for i := len(wordRunes) - 1; i >= 0; i-- {
		r := wordRunes[i]
		if !strings.Contains(excludeLetters, r) {
			last = r
			break
		}
	}
	return last
}

// nextWord устанавливает текущее слово и следующую букву
// func (g *Game) nextWord(word string) {
// 	g.word = strings.Title(word)
// 	g.nextRune = ""

// 	// разделить присланное слово на отдельные буквы
// 	wordRunes := strings.Split(strings.ToLower(word), "")

// 	// пройти по массиву букв начиная с последней и найти ту,
// 	// с которой будет начинаться следующее слово
// 	for i := len(wordRunes) - 1; i >= 0; i-- {
// 		if !strings.Contains(excludeLetters, wordRunes[i]) {
// 			g.nextRune = wordRunes[i]
// 			break
// 		}
// 	}
// }

// findNextPlayer ищет index следующего не выбывшего игрока, если не найден возвращает ok = false
func (g *Game) findNextPlayer() int {
	// оставшиеся игроки до конца массива
	for i, p := range g.players[g.nextPlayer+1:] {
		if !p.IsOut {
			return i
		}
	}

	// игроки с начала массива до текущего
	for i, p := range g.players[:g.nextPlayer] {
		if !p.IsOut {
			return i
		}
	}

	return 0
}

// activePlayers возвращает количество не выбывших игроков
func (g *Game) activePlayers() int {
	count := 0
	for _, p := range g.players {
		if !p.IsOut {
			count++
		}
	}
	return count
}

// nextTurn проверяет есть ли победитель и если есть рассылает уведомление о окончании матча
// иначе устанавливает следующего игрока и отсылает команду с запрсом нового слова
func (g *Game) nextTurn() {
	// есть победитель или не осталось игроков
	if g.activePlayers() <= 1 {
		panic("We have a winner")
	}

	// найти следующего игрока и разослать пакет
	g.nextPlayer = g.findNextPlayer()
	g.turnChan <- Turn{
		Word:    g.word,
		Letter:  g.nextRune,
		Timeout: g.turnTimeout,
		Player:  g.players[g.nextPlayer].Name,
	}
}

// RegisterPlayer добавляет игрока, игрок может быть добавлен только до начала партии.
func (g *Game) RegisterPlayer(name string, password string) error {
	if g.isRegistered(name) {
		return AlreadyRegistered
	}

	p := &Player{
		Name: name,
	}
	g.players = append(g.players, p)

	return nil
}

// MakeTurn принимает ответ от игрока, если такого игрока нет в списке возвращается ошибка UnknownPlayer.
// Пустой word считается за проигрыш.
// Не совпадение последней буквы и первой буквы присланного слова считается за проигрыш.
// Если слово отсутствует в словаре игрок выбывает.
// Если слово повторяется, то игрок выбывает.
// Если во время хода игрок выбывает, то ход переходит следующему игроку.
func (g *Game) MakeTurn(name string, word string) error {
	p, i, ok := g.findPlayer(name)
	if !ok {
		return UnknownPlayer
	}

	// ход другого игрока
	if i != g.nextPlayer {
		return WrongPlayer
	}

	// пустой ответ
	if utf8.RuneCountInString(word) == 0 {
		p.IsOut = true
		g.nextTurn()
		return EmptyWord
	}

	// неизвестное слово
	if !g.dictonary.HasKey(word) {
		p.IsOut = true
		g.nextTurn()
		return MissedInDictonary
	}

	// повтор
	if g.used.HasKey(word) {
		p.IsOut = true
		g.nextTurn()
		return UsedWord
	}

	// неверное слово
	if !g.validWord(word) {
		p.IsOut = true
		g.nextTurn()
		return Mistmatch
	}

	// с этого момента считается что игрок прислал верное слово
	g.used.SetKey(word)
	g.word = word
	g.nextRune = g.lastRune(word)
	g.nextTurn()
	return nil
}

// Start начинает игру
func (g *Game) Start() (<-chan Turn, error) {
	if g.started {
		return nil, AlredyStarted
	}

	g.used.Reset()
	g.turnChan = make(chan Turn, 1000)
	g.started = true
	return g.turnChan, nil
}
