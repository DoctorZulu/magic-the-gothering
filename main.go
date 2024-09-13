package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type Card struct {
	Name      string `json:"name"`
	ManaCost  string `json:"mana_cost"`
	Type      string `json:"type_line"`
	Power     string `json:"power"`
	Toughness string `json:"toughness"`
}

type Player struct {
	Name   string
	Hand   []Card
	Health int
	mu     sync.Mutex // Mutex to synchronize concurrent access to hand
}

// Constructor function to create a Player with default values
func NewPlayer(name string) *Player {
	p := &Player{
		Name:   name,
		Hand:   []Card{},
		Health: 100,
	}

	initHand(p)

	return p
}

func initHand(p *Player) {
	var wg sync.WaitGroup
	numRequests := 5

	// Queries 5 times and uses Mutex to avoid race condition updating Hand
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Fetch a random card
			response, err := http.Get("https://api.scryfall.com/cards/random")
			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}

			cards, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}

			var result Card
			err = json.Unmarshal(cards, &result)
			if err != nil {
				log.Fatalf("Failed to parse the JSON: %v", err)
			}

			p.mu.Lock()
			p.Hand = append(p.Hand, result)
			p.mu.Unlock()

			//fmt.Printf("Card Added: %+v\n", result)
		}()
	}
	// Wait for all goroutines to finish
	wg.Wait()
}

func main() {
	p1 := NewPlayer("Dan")
	p2 := NewPlayer("Seth")
	fmt.Printf("%+v\n", p1)
	fmt.Printf("%+v\n", p2)

}
