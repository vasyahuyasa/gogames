package main

import (
	"github.com/vasyahuyasa/gogames/words/core"
	udpTransport "github.com/vasyahuyasa/gogames/words/transport/udp"
	"log"
	"time"
)

const port = 0xbeef

func main() {
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
	turns, err := game.Start()
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

		// игра запросилас слово
		case turn := <-turns:
			if err := trans.SendTurn(turn); err != nil {
				trans.Error(turn.Player, err)
				log.Printf("Ошибка хода игрока %s: %v", turn.Player, err)
			}
		}
	}
}
