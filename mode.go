package main

var (
	modeOTT = map[uint32]button{
		0xFB04E786: PWR,
		0xFF00E786: MODE,
		0xBC43E786: AOUT,
		0xE817E786: UP,
		0xE41BE786: DOWN,
		0xE51AE786: LEFT,
		0xA659E786: RIGHT,
		0xBD42E786: OK,
		0xE916E786: BACK,
		0xF00FE786: HOME,
	}
	modeTV = map[uint32]button{
		0xFB04E786: PWR,
		0xFF00E786: MODE,
		0xBC43E786: INPUT,
		0xE817E786: CHUP,
		0xE41BE786: CHDOWN,
		0xE51AE786: VOLDOWN,
		0xA659E786: VOLUP,
		0xBD42E786: OK,
		0xE916E786: BACK,
		0xF00FE786: INFO,
	}

	modes = []map[uint32]button{
		modeOTT,
		modeTV,
	}
	currMode = 0
)
