package main

import (
	svr "gin-project/internal/server"
)

func main() {
	svr.FatalLog.PrintFatal(svr.NewServer(), nil)
}
