package core

import (
	"testing"
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

		g.nextPlayer = next
	}
}

func Test_activePlayers(t *testing.T) {
	g := gameForTest()
	active := g.activePlayers()
	if active != 4 {
		t.Fatalf("Активные игроки: ожидание 4, результат: %d", active)
	}
}
