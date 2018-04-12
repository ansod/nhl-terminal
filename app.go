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

var flags = map[string]bool{
	"goal":  true,
	"state": true,
}

func main() {

	if len(os.Args) > 1 {
		parser := argparse.NewParser(name, desc)

		goals := parser.Flag("g", "goal", &argparse.Options{Help: "Get alerted when a goal has been scored."})
		state := parser.Flag("s", "state", &argparse.Options{Help: "Get alerted when a game starts or ends."})

		err := parser.Parse(os.Args)

		if err != nil {
			// handle error
			fmt.Println(err)
		}

		flags["goal"] = *goals
		flags["state"] = *state
	}

	initGame()

	update()

}

func initGame() {
	response := helpers.Get()

	printScores(response)

	for _, game := range response.Dates[0].Games {
		gameScore[game.GamePK] = []int{game.Teams.Home.Score, game.Teams.Away.Score}
		gameState[game.GamePK] = game.Status.AbstractGameState
	}
}

func update() {
	for {
		msg, change := getUpdateMessage()

		if change {
			fmt.Println("==== UPDATE ====\n", msg)
		}

		time.Sleep(3 * time.Minute)
	}
}

func getUpdateMessage() (string, bool) {
	var updateMessage bytes.Buffer

	response := helpers.Get()

	for _, game := range response.Dates[0].Games {
		if flags["goal"] {
			s, change := compareScore(game)
			if change {
				updateMessage.WriteString(s)
			}
		}
		if flags["state"] {
			s, change := compareState(game)
			if change {
				updateMessage.WriteString(s)
			}
		}
	}

	return updateMessage.String(), updateMessage.Len() > 0
}

func printScores(response helpers.JSON) {
	var message bytes.Buffer

	for _, game := range response.Dates[0].Games {
		s := fmt.Sprintf("\n%s %d : %s %d\n",
			game.Teams.Home.Team.Abbreviation, game.Teams.Home.Score,
			game.Teams.Away.Team.Abbreviation, game.Teams.Away.Score)
		message.WriteString(s)
	}

	fmt.Println(message.String())
}

func compareScore(game helpers.Game) (string, bool) {
	hScore := game.Teams.Home.Score
	aScore := game.Teams.Away.Score
	hOldScore := gameScore[game.GamePK][0]
	aOldScore := gameScore[game.GamePK][1]

	if hScore == hOldScore && aScore == aOldScore {
		return "", false
	}

	gameScore[game.GamePK] = []int{hScore, aScore}

	if hScore > hOldScore && aScore > aOldScore {
		return fmt.Sprintf("\n%s [%d] : %s [%d]\n",
			game.Teams.Home.Team.Abbreviation, hScore,
			game.Teams.Away.Team.Abbreviation, aScore), true
	}

	if hScore > hOldScore {
		return fmt.Sprintf("\n%s [%d] : %s %d\n",
			game.Teams.Home.Team.Abbreviation, hScore,
			game.Teams.Away.Team.Abbreviation, aScore), true
	}

	return fmt.Sprintf("\n%s %d : %s [%d]\n",
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
