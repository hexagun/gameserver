package gameserver

type DisconnectPlayerPayload struct {
	Player string
}

type DisconnectPlayerAction struct {
	Type    string
	Payload DisconnectPlayerPayload
}

func (r DisconnectPlayerAction) GetType() string {
	return r.Type
}

func (r DisconnectPlayerAction) GetPayload() interface{} {
	return r.Payload
}

func CreateDisconnectPlayerAction(player string) DisconnectPlayerAction {
	return DisconnectPlayerAction{
		Type: "Disconnect",
		Payload: DisconnectPlayerPayload{
			Player: player,
		},
	}
}
