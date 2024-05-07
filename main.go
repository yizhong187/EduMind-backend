package main

import (
	"log"

	"github.com/yizhong187/EduMind-backend/cmd/api"
)

func main() {
	server := api.NewAPIServer(":8080", nil)
	if err := server.Run(); err != nil {
		log.Fatal("failed to run")
	}

}
