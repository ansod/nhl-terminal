package helpers

type JSON struct {
	Dates []Date
}

type Date struct {
	Games []Game `json:"games"`
}

// TODO: Add scoringplays
type Game struct {
	Status GameStatus `json:"status"`
	Teams  GameTeams
}

type GameStatus struct {
	AbstractGameState string
}

type GameTeams struct {
	Away AwayTeam
	Home HomeTeam
}

type AwayTeam struct {
	Score int
	Team  TeamInfo
}

type HomeTeam struct {
	Score int
	Team  TeamInfo
}

type TeamInfo struct {
	Name string
}
