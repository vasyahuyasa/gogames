package core

import (
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/vasyahuyasa/gogames/words/core/dictonary"
)

// Буквы на которые не может начинаться следующее слово
const excludeLetters = "цфзшэъьый"

// Game отвечает за всю игровую логику
type Game struct {
	players   []*Player
	dictonary *dictonary.Сollection
	used      *dictonary.Сollection

	started     bool
	turnTimeout time.Duration
	turnChan    chan Turn

	word          string
	nextRune      string
	currentPlayer int
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

// findNextPlayer ищет index следующего не выбывшего игрока
func (g *Game) findNextPlayer() int {
	// оставшиеся игроки до конца массива
	nextIndex := g.currentPlayer + 1
	for i, p := range g.players[nextIndex:] {
		if !p.IsOut() {
			return nextIndex + i
		}
	}

	// игроки с начала массива до текущего
	for i, p := range g.players[:g.currentPlayer] {
		if !p.IsOut() {
			return i
		}
	}

	return 0
}

// activePlayers возвращает количество не выбывших игроков
func (g *Game) activePlayers() int {
	count := 0
	for _, p := range g.players {
		if !p.IsOut() {
			count++
		}
	}
	return count
}

// sendTurn отсылает команду с запрсом нового слова текущему игроку
func (g *Game) sendTurn() {
	// есть победитель или не осталось игроков
	if g.activePlayers() <= 1 {
		panic("We have a winner")
	}

	log.Printf("отправка запроса слова игроку: %q", g.players[g.currentPlayer].Name)

	// найти следующего игрока и разослать пакет
	g.turnChan <- Turn{
		Word:    g.word,
		Letter:  g.nextRune,
		Timeout: g.turnTimeout,
		Player:  g.players[g.currentPlayer].Name,
	}
}

// RegisterPlayer добавляет игрока, игрок может быть добавлен только до начала партии.
func (g *Game) RegisterPlayer(name string) error {
	if g.started {
		return GameInProgress
	}

	if g.isRegistered(name) {
		return AlreadyRegistered
	}

	p := &Player{
		Name: name,
	}
	g.players = append(g.players, p)

	log.Printf("Игрок зарегестрирован: %q", p.Name)

	return nil
}

// checkTurn проверяет является ли слово word подходящим для ответа.
// за неверный ответ считается:
// Пустой word;
// Не совпадение последней буквы и первой буквы присланного слова;
// Если слово отсутствует в словаре;
// Если слово повторяется;
func (g *Game) checkWord(word string) error {
	// пустой ответ
	if utf8.RuneCountInString(word) == 0 {
		return EmptyWord
	}

	// неизвестное слово
	if !g.dictonary.HasKey(word) {
		return MissedInDictonary
	}

	// повтор
	if g.used.HasKey(word) {
		return UsedWord
	}

	// неверное слово
	if !g.validWord(word) {
		return Mistmatch
	}

	return nil
}

// MakeTurn принимает ответ от игрока, если такого игрока нет в списке возвращается ошибка UnknownPlayer.
// Если во время хода игрок выбывает, то ход переходит следующему игроку.
func (g *Game) MakeTurn(name string, word string) error {
	p, i, ok := g.findPlayer(name)
	if !ok {
		return UnknownPlayer
	}

	// ход другого игрока
	if i != g.currentPlayer {
		return WrongPlayer
	}

	// с этого момента ход переходит другому игроку
	g.currentPlayer = g.findNextPlayer()

	// верное ли слово прислал игрок
	err := g.checkWord(word)
	if err != nil {
		p.Out()
		g.sendTurn()
		return err
	}

	// с этого момента считается что игрок прислал верное слово
	g.used.SetKey(word)
	g.word = word
	g.nextRune = g.lastRune(word)
	g.sendTurn()
	return nil
}

// Start начинает игру
func (g *Game) Start(dic *dictonary.Сollection) (<-chan Turn, error) {
	if dic == nil {
		return nil, fmt.Errorf("пустой словарь")
	}

	if g.started {
		return nil, AlredyStarted
	}

	if len(g.players) == 0 {
		return nil, NotEnoughtPlayers
	}

	g.dictonary = dic
	g.used = dictonary.New()
	g.turnChan = make(chan Turn, 1000)
	g.started = true

	// игрок с индексом 0 в игре
	g.sendTurn()

	return g.turnChan, nil
}
