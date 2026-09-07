package main

import (
	"nextcmd/cmd/assistant"
	"nextcmd/internal/windowicon"
)

func main() {
	windowicon.Apply()
	assistant.Main()
}
