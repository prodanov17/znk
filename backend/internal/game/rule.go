package game

func DefaultRules() map[string]int {
	m := make(map[string]int)

	m["winning_score"] = 120
	m["see_dream_card"] = 1
	m["show_capture"] = 1

	return m
}
