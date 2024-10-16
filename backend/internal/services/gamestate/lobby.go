package gamestate

type Lobby struct {
	ID        string `json:"id"`
	OwnerID   int    `json:"owner_id"`
	GameID    string `json:"game_id"`
	CreatedAt string `json:"created_at"`
	GameTeam  []GameTeam
}
