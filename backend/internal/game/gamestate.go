package game

import "fmt"

type GameState struct {
	GameRound int `json:"game_round"`
	TurnIdx   int `json:"turn_idx"`
	DealerIdx int `json:"dealer_idx"`
}

func NewGameState() *GameState {
	return &GameState{
		GameRound: 1,
		TurnIdx:   1,
		DealerIdx: 0,
	}
}

func (gs *GameState) NextTurn(teams []GameTeam) (*Player, error) {
	if len(teams) != 2 || len(teams[0].Players) != 2 || len(teams[1].Players) != 2 {
		return nil, fmt.Errorf("invalid team or player setup, must be 2 teams of 2 players each")
	}

	turnIdx := gs.TurnIdx % 4

	switch turnIdx {
	case 0:
		return teams[0].Players[0], nil // Team 1, Player 1
	case 1:
		return teams[1].Players[0], nil // Team 2, Player 1
	case 2:
		return teams[0].Players[1], nil // Team 1, Player 2
	case 3:
		return teams[1].Players[1], nil // Team 2, Player 2
	}

	return nil, fmt.Errorf("unexpected turn index")
}

func (gs *GameState) AdvanceTurn() {
	gs.TurnIdx++
	if (gs.TurnIdx-(gs.DealerIdx+1))%52 == 0 { // all cards have been played
		gs.GameRound++
	}
}

func (gs *GameState) Dealer(teams []GameTeam) (*Player, error) {
	if len(teams) != 2 || len(teams[0].Players) != 2 || len(teams[1].Players) != 2 {
		return nil, fmt.Errorf("invalid team or player setup, must be 2 teams of 2 players each")
	}

	dealerIdx := gs.DealerIdx % 4

	switch dealerIdx {
	case 0:
		return teams[0].Players[0], nil // Team 1, Player 1
	case 1:
		return teams[1].Players[0], nil // Team 2, Player 1
	case 2:
		return teams[0].Players[1], nil // Team 1, Player 2
	case 3:
		return teams[1].Players[1], nil // Team 2, Player 2
	}

	return nil, fmt.Errorf("unexpected dealer index")

}

func (gs *GameState) AdvanceDealer() {
	gs.DealerIdx++
}

func (gs *GameState) IsRoundOver() bool {
	return (gs.TurnIdx-(gs.DealerIdx+1))%52 == 0 // since turnidx starts at dealeridx+1
}

func (gs *GameState) Reset() {
	gs.GameRound = 1
	gs.TurnIdx = 1
	gs.DealerIdx = 0
}
