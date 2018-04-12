package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/akamensky/argparse"
	"github.com/ansod/nhl-terminal/helpers"
)

const (
	name = "nhl-terminal"
	desc = "" // TODO: Add description
)

// Contains the score from last check.
// {gameID: {homeScore, AwayScore}}
var gameScore = map[int][]int{}

// Contains the game state from last check.
// {gameID: state}
var gameState = map[int]string{}

func main() {

	var flags = map[string]bool{}

	parser := argparse.NewParser(name, desc)

	goals := parser.Flag("g", "goal", &argparse.Options{Help: "Get alerted when a goal has been scored."})
	start := parser.Flag("s", "start", &argparse.Options{Help: "Get alerted when a game starts."})
	end := parser.Flag("e", "end", &argparse.Options{Help: "Get alerted when a game has ended."})

	flags["goal"] = *goals
	flags["start"] = *start
	flags["end"] = *end

	err := parser.Parse(os.Args)

	if err != nil {
		// handle error
		fmt.Println(err)
	}

	initGame()

	update(flags)

}

func initGame() {
	response := helpers.Get()

	for _, game := range response.Dates[0].Games {
		gameScore[game.GamePK] = []int{game.Teams.Home.Score, game.Teams.Away.Score}
		gameState[game.GamePK] = game.Status.AbstractGameState
	}
}

func update(flags map[string]bool) {
	var updateMessage bytes.Buffer

	for {
		fmt.Println("reading...")
		response := helpers.Get()

		for _, game := range response.Dates[0].Games {
			s, change := compareScore(game)
			if change {
				updateMessage.WriteString(s)
			}
			s, change = compareState(game)
			if change {
				updateMessage.WriteString(s)
			}
		}

		fmt.Println(updateMessage.String())

		updateMessage.Reset()

		time.Sleep(1 * time.Minute)
	}
}

func compareScore(game helpers.Game) (string, bool) {
	hScore := game.Teams.Home.Score
	aScore := game.Teams.Away.Score

	if hScore == gameScore[game.GamePK][0] && aScore == gameScore[game.GamePK][1] {
		return "", false
	}

	gameScore[game.GamePK] = []int{hScore, aScore}

	return fmt.Sprintf("Game score: %s %d : %s %d\n",
		game.Teams.Home.Team.Abbreviation, hScore,
		game.Teams.Away.Team.Abbreviation, aScore), true

}

func compareState(game helpers.Game) (string, bool) {
	state := game.Status.AbstractGameState

	if state == gameState[game.GamePK] {
		return "", false
	}

	gameState[game.GamePK] = state

	return fmt.Sprintf("Game state has changed from %s to %s\n", gameState[game.GamePK], state), true
}
