package gameserver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayTokenAction(t *testing.T) {

	var playerName = "Player1"
	var x = 1
	var y = 2
	var action Action
	action = CreatePlayTokenAction(playerName, x, y)
	var playTokenAction PlayTokenAction = action.(PlayTokenAction)
	var payload PlayTokenPayload = playTokenAction.GetPayload().(PlayTokenPayload)
	assert.Equal(t, "PlayToken", playTokenAction.GetType(), "wrong action type")
	assert.Equal(t, playerName, payload.Player, "Payload Player is not named Player1")
	assert.Equal(t, x, payload.X, "Payload tile X not equal to 1")
	assert.Equal(t, y, payload.Y, "Payload tile Y not equal to 2")
}

func TestConnectPlayerAction(t *testing.T) {

	var playerName = "Player1"

	var action Action
	action = CreateConnectPlayerAction(playerName)
	var connectPlayerAction ConnectPlayerAction = action.(ConnectPlayerAction)
	var payload ConnectPlayerPayload = connectPlayerAction.GetPayload().(ConnectPlayerPayload)
	assert.Equal(t, "Connect", connectPlayerAction.GetType(), "wrong action type")
	assert.Equal(t, playerName, payload.Player, "Payload Player is not named Player1")
}

func TestDisconnectPlayerAction(t *testing.T) {

	var playerName = "Player1"

	var action Action
	action = CreateDisconnectPlayerAction(playerName)
	var disconnectPlayerAction DisconnectPlayerAction = action.(DisconnectPlayerAction)
	var payload DisconnectPlayerPayload = disconnectPlayerAction.GetPayload().(DisconnectPlayerPayload)
	assert.Equal(t, "Disconnect", disconnectPlayerAction.GetType(), "wrong action type")
	assert.Equal(t, playerName, payload.Player, "Payload Player is not named Player1")
}
