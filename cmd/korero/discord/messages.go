package discord

import (
	"fmt"
	"korero/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

var (
	messagesDescription = `logs in to the default user and streams messages from their server`

	DiscordMessagesCommand = &cobra.Command{
		Use:     "messages TOKEN",
		Short:   "stream messages from the given user",
		Long:    messagesDescription,
		RunE:    messages,
		Example: "korero discord messages 12345",
	}
)

func messagesFlags() {

}

func init() {
	messagesFlags()
}

func messages(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		Token = args[0]
	}
	dg, err := utils.LoginDiscord(Token)
	if err != nil {
		return err
	}

	dg.AddHandler(list)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		return err
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()

	return nil
}

func list(dg *discordgo.Session, message *discordgo.MessageCreate) {
	fmt.Println(message.Content)
}
