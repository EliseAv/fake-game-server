package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	fmt.Println("Loading . . .")
	rand.Seed(time.Now().Unix())
	randSleep(500, 30)
	fmt.Println("Running!!!")

	go simulatePlayerActivity()
	go simulateAutosave()
	readStdin()
	fmt.Println("Closing server . . .")
	randSleep(500, 30)
	fmt.Println("Done.")
}

func randSleep(min time.Duration, variance int64) {
	x := time.Duration(rand.Int63n(variance))
	time.Sleep(time.Millisecond * (min + x*x))
}

func readStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "quit" {
			return
		}
		fmt.Println("<server>", text)
	}
}

const (
	joinChance            = 1
	talkChance1Player     = 1
	talkChanceManyPlayers = 8
	leaveChance           = 1
)

type player string

var greets = []string{"hey", "hi", "hello", "sup", "yo", "uh"}
var offlinePlayers = []player{
	"mike", "jess", "chris", "ash", "matt", "brit", "josh", "amy", "dan", "sam",
	"dave", "sarah", "andy", "steph", "jay", "jenn", "jus", "beth", "jo", "lauren",
}
var onlinePlayers = []player{}

func simulatePlayerActivity() {
	for {
		var p player
		talkChance := talkChanceManyPlayers
		randSleep(100, 100)
		currNumPlayers := len(onlinePlayers)
		if currNumPlayers == 0 {
			p.join()
			continue
		} else if currNumPlayers == 1 {
			talkChance = talkChance1Player
		}

		playerSpace := talkChance + leaveChance
		diceRoll := rand.Intn(currNumPlayers*playerSpace + joinChance)
		selectedPlayer := diceRoll / playerSpace
		if selectedPlayer >= currNumPlayers {
			p.join()
			continue
		}
		p = onlinePlayers[selectedPlayer]
		if (diceRoll % playerSpace) < talkChance {
			p.talk()
		} else {
			p.leave()
		}
	}
}

func (p player) join() {
	numBefore := len(offlinePlayers)
	if numBefore == 0 {
		return
	}
	iRand := rand.Intn(numBefore)
	p = offlinePlayers[iRand]
	iLast := numBefore - 1
	offlinePlayers[iRand] = offlinePlayers[iLast]
	offlinePlayers = offlinePlayers[:iLast]

	onlinePlayers = append(onlinePlayers, p)

	fmt.Printf("* %s joined!\n", p)
}

func (p player) talk() {
	num := len(onlinePlayers)
	if num == 1 {
		// Lone player
		fmt.Printf("<%s> hm%*s\n", p, 3+rand.Intn(4), "...")
		return
	}

	other := onlinePlayers[rand.Intn(num-1)]
	if other == p {
		other = onlinePlayers[num-1]
	}

	fmt.Printf("<%s> %s %s\n", p, greets[rand.Intn(len(greets))], other)
}

func (p player) leave() {
	offlinePlayers = append(offlinePlayers, p)

	slice := onlinePlayers[:0]
	for _, x := range onlinePlayers {
		if x != p {
			slice = append(slice, x)
		}
	}
	onlinePlayers = slice

	fmt.Printf("* %s left!\n", p)
}

func simulateAutosave() {
	tick := time.Tick(1 * time.Minute)
	for range tick {
		fmt.Println("Autosaving . . .")
		randSleep(100, 10)
		fmt.Println("Autosaved.")
	}
}
