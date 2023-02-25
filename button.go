package main

type button uint8

const (
	PWR button = iota
	MODE
	AOUT
	UP
	DOWN
	LEFT
	RIGHT
	OK
	BACK
	HOME
	INPUT
	CHUP
	CHDOWN
	VOLDOWN
	VOLUP
	CHLIST
)

func (b button) String() string {
	switch b {
	case PWR:
		return "PWR"
	case MODE:
		return "MODE"
	case AOUT:
		return "AOUT"
	case UP:
		return "UP"
	case DOWN:
		return "DOWN"
	case LEFT:
		return "LEFT"
	case RIGHT:
		return "RIGHT"
	case OK:
		return "OK"
	case BACK:
		return "BACK"
	case HOME:
		return "HOME"
	case INPUT:
		return "INPUT"
	case CHUP:
		return "CHUP"
	case CHDOWN:
		return "CHDOWN"
	case VOLDOWN:
		return "VOLDOWN"
	case VOLUP:
		return "VOLUP"
	case CHLIST:
		return "CHLIST"
	default:
		return "UNKNOWN"
	}
}
