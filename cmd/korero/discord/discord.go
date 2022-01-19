package discord

import "github.com/spf13/cobra"

var (
	DiscordCommand = &cobra.Command{
		Use:              "discord",
		Short:            "Connect to and use the discord API",
		Long:             "Connect to and use the discord API",
		TraverseChildren: true,
	}
)

func init() {
	DiscordCommand.AddCommand(DiscordLoginCommand)
	DiscordCommand.AddCommand(DiscordMessagesCommand)
	DiscordCommand.AddCommand(DiscordSetupCommand)
}
