package main

import (
	"gin-project/internal/server"
)

func main() {

	logger := server.InfoLog
	server := server.NewServer()

	err := server.ListenAndServe()

	logger.PrintFatal(err, nil)
}
