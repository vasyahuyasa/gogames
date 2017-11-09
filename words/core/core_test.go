package core

import (
	"github.com/vasyahuyasa/gogames/words/core/dictonary"
	"testing"
	"time"
)

func gameForTest() *Game {
	return &Game{
		players: []*Player{
			&Player{"Test", true},      // 0
			&Player{"Петр", false},     // 1
			&Player{"Иванович", false}, // 2
			&Player{"ж", false},        // 3
			&Player{"Outsider", true},  // 4
			&Player{"Lost", true},      // 5
			&Player{"Active", false},   // 6
		},
	}
}

func testCollection() *dictonary.Сollection {
	dict := []string{
		"test",
		"test2",
		"word",
		"hi",
	}
	return dictonary.New(dict...)
}

func TestGame_findPlayer(t *testing.T) {
	g := gameForTest()

	tests := []struct {
		name  string
		index int
		ok    bool
	}{
		{"Test", 0, true},
		{"Петр", 1, true},
		{"ж", 3, true},
		{"no", 0, false},
		{"Пет", 0, false},
	}

	for _, test := range tests {
		p, index, ok := g.findPlayer(test.name)
		if ok != test.ok {
			t.Fatalf("Поиск: ожидание %t, результат %t", test.ok, ok)
		}

		if ok {
			if p.Name != test.name {
				t.Fatalf("Имя: ожидание %q, результат %q", test.name, p.Name)
			}

			if index != test.index {
				t.Fatalf("Индекс: ожидание %q, результат %q", test.index, index)
			}
		}
	}
}

func TestGame_lastRune(t *testing.T) {
	tests := map[string]string{
		"тестц":       "т",
		"сызрань":     "н",
		"Вологда":     "а",
		"Адыгея":      "я",
		"Ялта":        "а",
		"Архангельск": "к",
		"Test": "t",
	}

	g := gameForTest()

	for word, last := range tests {
		l := g.lastRune(word)
		if l != last {
			t.Fatalf("Неверная результат: ожидание %q, результат %q", last, l)
		}
	}
}

func Test_findNextPlayer(t *testing.T) {
	tests := []int{1, 2, 3, 6, 1, 2, 3, 6}
	g := gameForTest()

	for _, test := range tests {
		next := g.findNextPlayer()

		if next != test {
			t.Fatalf("Найден неверный игрок: ожидание %d, результат %d", test, next)
		}

		g.currentPlayer = next
	}
}

func Test_activePlayers(t *testing.T) {
	g := gameForTest()
	active := g.activePlayers()
	if active != 4 {
		t.Fatalf("Активные игроки: ожидание 4, результат: %d", active)
	}
}

func Test_sendTurn(t *testing.T) {
	test := struct {
		word     string
		nextRune string
		timeout  time.Duration
		player   int
	}{
		word:     "test",
		nextRune: "x",
		timeout:  100,
		player:   2,
	}

	g := gameForTest()

	g.turnChan = make(chan Turn, 1)
	g.word = test.word
	g.nextRune = test.nextRune
	g.turnTimeout = test.timeout
	g.currentPlayer = test.player

	g.sendTurn()
	r := <-g.turnChan

	if r.Word != test.word ||
		r.Letter != test.nextRune ||
		r.Timeout != test.timeout ||
		r.Player != g.players[test.player].Name {
		t.Fatal("Одно или несколько полей не совпадают с внутренним состоянием")
	}
}

func Test_RegisterPlayer(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "player1",
			err:  nil,
		},
		{
			name: "player2",
			err:  nil,
		},
		{
			name: "player1",
			err:  AlreadyRegistered,
		},
	}

	g := gameForTest()

	for _, test := range tests {
		err := g.RegisterPlayer(test.name)

		// добавление игрока
		if err != test.err {
			t.Fatalf("Регистрация игрока %s: ожидание %v результат %v", test.name, test.err, err)
		}

		// проверить есть ли в списке
		if err == nil {
			p := g.players[len(g.players)-1]
			if p.Name != test.name {
				t.Fatal("Игрок не добавлен")
			}
		}
	}

	// добавление во время запущенной игры
	expected := GameInProgress
	g.started = true
	err := g.RegisterPlayer("random888")
	if err != expected {
		t.Fatalf("Добавление игрока во время запущенной игры: ожидание %v, результат %v", expected, err)
	}
}
