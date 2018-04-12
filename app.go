package main

import (
	"fmt"
	"os"

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

	update(flags)

}

func update(flags map[string]bool) {
	response := helpers.Get()

	fmt.Println(response.Dates[0].Games[0].Status.AbstractGameState)
}

func compareScore(game helpers.Game) (string, bool) {
	hScore := game.Teams.Home.Score
	aScore := game.Teams.Away.Score

	if hScore == gameScore[game.GamePK][0] && aScore > gameScore[game.GamePK][1] {
		return "", false
	}

	return fmt.Sprintf("Game score: %s %d : %s %d",
		game.Teams.Home.Team.Abbreviation, hScore,
		game.Teams.Away.Team.Abbreviation, aScore), true

}

func compareState(game helpers.Game) (string, bool) {
	state := game.Status.AbstractGameState

	if state == gameState[game.GamePK] {
		return "", false
	}

	return fmt.Sprintf("Game state has changed from %s to %s", gameState[game.GamePK], state), true
}
