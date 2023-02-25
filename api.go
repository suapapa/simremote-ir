package main

import "github.com/go-resty/resty/v2"

var (
	apiClient = resty.New()
)

func handle(btn button) error {
	switch btn {
	case UP:
		_, err := apiClient.R().Put(apiAddr + "/key/up")
		return err
	case DOWN:
		_, err := apiClient.R().Put(apiAddr + "/key/down")
		return err
	case LEFT:
		_, err := apiClient.R().Put(apiAddr + "/key/left")
		return err
	case RIGHT:
		_, err := apiClient.R().Put(apiAddr + "/key/right")
		return err
	case OK:
		_, err := apiClient.R().Put(apiAddr + "/key/ok")
		return err
	case BACK:
		_, err := apiClient.R().Put(apiAddr + "/key/back")
		return err
	case HOME:
		_, err := apiClient.R().Put(apiAddr + "/key/home")
		return err
	case CHUP:
		_, err := apiClient.R().Put(apiAddr + "/key/chup")
		return err
	case CHDOWN:
		_, err := apiClient.R().Put(apiAddr + "/key/chdown")
		return err
	case VOLDOWN:
		_, err := apiClient.R().Put(apiAddr + "/key/voldown")
		return err
	case VOLUP:
		_, err := apiClient.R().Put(apiAddr + "/key/volup")
		return err
	}
	return nil
}
