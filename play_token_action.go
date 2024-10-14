package gameserver

type PlayTokenPayload struct {
	Player string
	X, Y   int
}

type PlayTokenAction struct {
	Type    string
	Payload PlayTokenPayload
}

func (r PlayTokenAction) GetType() string {
	return r.Type
}

func (r PlayTokenAction) GetPayload() interface{} {
	return r.Payload
}

func CreatePlayTokenAction(player string, x int, y int) PlayTokenAction {
	return PlayTokenAction{
		Type: "PlayToken",
		Payload: PlayTokenPayload{
			Player: player,
			X:      x,
			Y:      y,
		},
	}
}
