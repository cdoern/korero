package main

import (
	"korero/cmd/korero/discord"
	"os"
)

func main() {

	parent := rootCmd

	parent.AddCommand(discord.DiscordCommand)

	Execute()
	os.Exit(0)
}
