package main

import (
	"github.com/lucaspopp0/hass-update-manager/update-manager/api"
)

func main() {
	server := api.NewServer()
	server.Run()
}
