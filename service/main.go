package main

import (
	"github.com/haswelliris/service-agenda/service/service"
)

const (
	port string = "8080"
)

func main() {
	server := service.NewServer()
	server.Run(":" + port)

}
