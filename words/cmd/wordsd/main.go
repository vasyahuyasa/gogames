//go:generate go-bindata -nometadata -o cities_data.go cities.txt
package main

import (
	"bytes"
	"log"
	"time"

	"github.com/vasyahuyasa/gogames/words/core"
	"github.com/vasyahuyasa/gogames/words/core/dictonary"
	udpTransport "github.com/vasyahuyasa/gogames/words/transport/udp"
)

const port = 0xbeef

func main() {
	data, err := Asset("cities.txt")
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(data)
	dict := dictonary.New()
	dict.FromReader(r)
	log.Print("Словарь загружен")

	game := &core.Game{}
	trans, err := udpTransport.New(port)
	if err != nil {
		log.Fatalf("Инициализация udp транспорта: %v", err)
	}

	// ожидание регистраций
	go func() {
		for reg := range trans.RegChan() {
			if err := game.RegisterPlayer(reg.Name); err != nil {
				trans.Error("", err)
				log.Printf("Ошибка регистрации игрока в игре: %v", err)
			}
		}
	}()

	log.Println("Начало игры через 5 сек")
	time.Sleep(time.Second * 5)
	turns, err := game.Start(dict)
	if err != nil {
		log.Fatalf("Невозможно начать игру: %v", err)
	}
	log.Println("Игра началась...")

	for {
		select {
		// игрок прислал ответ
		case turn := <-trans.TurnChan():
			if err := game.MakeTurn(turn.Name, turn.Word); err != nil {
				trans.Error("", err)
				log.Printf("Ошибка регистрации игрока в игре: %v", err)
			}

		// игра запросила слово
		case turn := <-turns:
			if err := trans.SendTurn(turn); err != nil {
				trans.Error(turn.Player, err)
				log.Printf("Ошибка хода игрока %s: %v", turn.Player, err)
			}
		}
	}
}
