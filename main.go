package main

import (
	"tictactoe/internal/di"
)

func main() {
	app := di.NewApp()
	app.Run()
}
