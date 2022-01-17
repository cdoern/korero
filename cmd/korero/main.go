package main

import (
	"korero/cmd/korero/discord"
	"os"
)

func main() {

	parent := rootCmd

	parent.AddCommand(discord.DiscordCommand)
	//parent.AddCommand(discord.DiscordLoginCommand)

	Execute()
	os.Exit(0)
}
