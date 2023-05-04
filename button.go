package main

type button uint8

const (
	UNKNOWN button = iota
	PWRON
	PWROFF
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
	INFO
)

func (b button) String() string {
	switch b {
	case UNKNOWN:
		return "UNKNOWN"
	case PWRON:
		return "PWRON"
	case PWROFF:
		return "PWROFF"
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
	case INFO:
		return "INFO"
	default:
		return "UNKNOWN"
	}
}
