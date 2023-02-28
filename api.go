package main

import (
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

type apiClient struct {
	baseAddr string
	tv       *TV
	client   *http.Client
}

func NewAPIClient(baseAddr string) *apiClient {
	return &apiClient{
		baseAddr: baseAddr,
		tv:       NewTV(),
		client: &http.Client{
			Timeout: 1 * time.Second,
		},
	}
}

func (c *apiClient) put(path string) error {
	req, err := http.NewRequest(http.MethodPut, c.baseAddr+path, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

/*
func (c *apiClient) get(path string) (map[string]any, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseAddr+path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return data["data"].(map[string]any), nil
}
*/

func (c *apiClient) Handle(btn button) error {
	switch btn {
	case UP, DOWN, LEFT, RIGHT, OK, BACK, HOME, CHUP, CHDOWN, VOLDOWN, VOLUP, INFO:
		return c.keyHandler(btn)
	case PWR:
		if c.tv.Status == tvStatusOff {
			if err := c.put("/tv/on"); err != nil {
				return errors.Wrap(err, "failed to turn on TV")
			}
			c.tv.Status = tvStatusOn
		} else {
			// if current status is on or unknown
			if err := c.put("/tv/off"); err != nil {
				if os.IsTimeout(err) {
					// if timeout, assume TV is already off
					c.tv.Status = tvStatusOff
					return nil
				}
				return errors.Wrap(err, "failed to turn off TV")
			}
			c.tv.Status = tvStatusOff
		}
	case AOUT:
		if err := c.put("/audio/" + c.tv.AudioOuts[c.tv.CurAppIdx]); err != nil {
			return errors.Wrap(err, "failed to change audio output")
		}
		c.tv.CurAppIdx = (c.tv.CurAppIdx + 1) % len(c.tv.AudioOuts)
	case INPUT:
		if err := c.put("/app/" + c.tv.Apps[c.tv.CurAppIdx]); err != nil {
			return errors.Wrap(err, "failed to change app")
		}
		c.tv.CurAppIdx = (c.tv.CurAppIdx + 1) % len(c.tv.Apps)
	}
	return nil
}

func (c *apiClient) keyHandler(btn button) error {
	btn2Key := map[button]string{
		UP:      "up",
		DOWN:    "down",
		LEFT:    "left",
		RIGHT:   "right",
		OK:      "ok",
		BACK:    "back",
		HOME:    "home",
		CHUP:    "channel_up",
		CHDOWN:  "channel_down",
		VOLDOWN: "volume_down",
		VOLUP:   "volume_up",
		INFO:    "info",
	}

	key, ok := btn2Key[btn]
	if !ok {
		return errors.Errorf("unknown button: %s", btn)
	}

	return c.put("/key/" + key)
}
