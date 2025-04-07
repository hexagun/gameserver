package main

type State int

// States
// No one is connected, game not started      		// InitialState
// One is connected, game not started				// WaitingForOpponent
// Two is connected, game not started				// ReadyToStart

// Two is connected, game started, ongoing			// InGame
// One is connected, game started, ongoing			// Rejoin
// No one is connected, game started, ongoing		// DoubleDisconnect

// Two is connected, game ended 					// GameEnded
// One is connected, game ended 					// GameOver
// No one is connected, game ended 					// FinalState

const (
	InitialState State = iota
	WaitingForOpponent
	ReadyToStart
	InGame
	Rejoin
	DoubleDisconnect
	GameEnded
	GameOver
	FinalState
)

func (s State) String() string {
	switch s {
	case InitialState:
		return "InitialState"
	case WaitingForOpponent:
		return "WaitingForOpponent"
	case ReadyToStart:
		return "ReadyToStart"
	case InGame:
		return "InGame"
	case Rejoin:
		return "Rejoin"
	case DoubleDisconnect:
		return "DoubleDisconnect"
	case GameEnded:
		return "GameEnded"
	case GameOver:
		return "GameOver"
	case FinalState:
		return "FinalState"
	}
	return "Unknown"
}

type Trigger int

const (
	Join Trigger = iota
	Leave
	Start
	Move
	// CallDialed Trigger = iota
	// HungUp
	// CallConnected
	// PlacedOnHold
	// TakenOffHold
	// LeftMessage
)

func (t Trigger) String() string {
	switch t {
	case Join:
		return "Join"
	case Leave:
		return "Leave"
	case Start:
		return "Start"
	case Move:
		return "Move"
	}
	return "Error"
}

type TriggerResult struct {
	Trigger Trigger
	State   State
}

var rules = map[State][]TriggerResult{
	InitialState: {
		{Join, WaitingForOpponent},
	},
	WaitingForOpponent: {
		{Join, ReadyToStart},
		{Leave, InitialState},
	},
	ReadyToStart: {
		{Start, InGame},
		{Leave, WaitingForOpponent},
	},
	InGame: {
		{Leave, Rejoin},
		{Move, InGame},
		{Move, GameEnded},
	},
	Rejoin: {
		{Leave, FinalState},
		{Join, InGame},
	},
	GameEnded: {
		{Leave, GameEnded},
		{Leave, FinalState},
	},

	FinalState: {},
}

// DoubleDisconnect: {
// },
// GameOver: {
// },

// func main() {
// 	state, exitState := InitialState, FinalState
// 	for ok := true; ok; ok = state != exitState {
// 		fmt.Println("The phone is currently", state)
// 		fmt.Println("Select a trigger:")

// 		for i := 0; i < len(rules[state]); i++ {
// 			tr := rules[state][i]
// 			fmt.Println(strconv.Itoa(i), ".", tr.Trigger)
// 		}

// 		input, _, _ := bufio.NewReader(os.Stdin).ReadLine()
// 		i, _ := strconv.Atoi(string(input))

// 		tr := rules[state][i]
// 		state = tr.State
// 	}
// 	fmt.Println("We are done using the phone")
// }
