package game

type GameTeam struct {
	ID      int    `json:"team_id"`
	GameID  string `json:"game_id"`
	Score   int    `json:"score"`
	Players []*Player
}

func NewGameTeam(id int, gameID string) *GameTeam {
	return &GameTeam{
		ID:     id,
		GameID: gameID,
		Score:  0,
	}
}

func (gt *GameTeam) AddPlayer(p *Player) {
	gt.Players = append(gt.Players, p)
}

func (gt *GameTeam) RemovePlayer(p *Player) {
	for i, player := range gt.Players {
		if player.UserID == p.UserID {
			gt.Players = append(gt.Players[:i], gt.Players[i+1:]...)
			break
		}
	}
}
