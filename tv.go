package main

type tvStatus uint8

const (
	tvStatusUnknown tvStatus = iota
	tvStatusOn
	tvStatusOff
)

type TV struct {
	Status         tvStatus
	AudioOuts      []string
	CurAudioOutIdx int
	Apps           []string
	CurAppIdx      int
}

func NewTV() *TV {
	return &TV{
		// TODO: get list from TV
		AudioOuts: []string{
			"tv_speaker",
			"external_optical",
		},
		Apps: []string{
			"현재 방송",
			"HDMI1",
			"HDMI2",
			// "HDMI3",
			// "HDMI4",
		},
	}
}
