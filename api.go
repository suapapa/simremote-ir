package main

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type apiClient struct {
	baseAddr string
}

func newAPIClient(baseAddr string) *apiClient {
	return &apiClient{
		baseAddr: baseAddr,
	}
}

func (c *apiClient) put(path string) error {
	req, err := http.NewRequest(http.MethodPut, c.baseAddr+path, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

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

func (c *apiClient) Handle(btn button) error {
	switch btn {
	case UP, DOWN, LEFT, RIGHT, OK, BACK, HOME, CHUP, CHDOWN, VOLDOWN, VOLUP:
		return c.keyHandler(btn)
	case PWR:
	case MODE:
	case AOUT:
	case CHLIST:
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
		CHUP:    "chup",
		CHDOWN:  "chdown",
		VOLDOWN: "voldown",
		VOLUP:   "volup",
	}

	key, ok := btn2Key[btn]
	if !ok {
		return errors.Errorf("unknown button: %s", btn)
	}

	return c.put("/key/" + key)
}
