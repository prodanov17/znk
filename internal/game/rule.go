package game

type GameRule struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func NewGameRule(id int, name string, value int) *GameRule {
	return &GameRule{
		ID:    id,
		Name:  name,
		Value: value,
	}
}

func DefaultRules() []GameRule {
	return []GameRule{
		*NewGameRule(1, "play_until", 120),
		*NewGameRule(2, "see_dream_card", 1),
	}
}
