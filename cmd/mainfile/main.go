package main

import (
	"DevIntApp/internal/api"
	"log"
)

func main() {
	log.Println("App started!")
	api.StartServer()
	log.Println("App finished!")
}
