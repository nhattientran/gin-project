package main

import (
	svr "gin-project/internal/server"
)

func main() {
	server := svr.NewServer()
	err := server.ListenAndServe()

	svr.InfoLog.PrintFatal(err, nil)
}
