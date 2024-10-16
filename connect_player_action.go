package gameserver

type ConnectPlayerPayload struct {
	Player string
}

type ConnectPlayerAction struct {
	Type    string
	Payload ConnectPlayerPayload
}

func (r ConnectPlayerAction) GetType() string {
	return r.Type
}

func (r ConnectPlayerAction) GetPayload() interface{} {
	return r.Payload
}

func CreateConnectPlayerAction(player string) ConnectPlayerAction {
	return ConnectPlayerAction{
		Type: "ConnectPlayer",
		Payload: ConnectPlayerPayload{
			Player: player,
		},
	}
}
