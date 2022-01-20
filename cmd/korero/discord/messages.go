package discord

import (
	"bufio"
	"fmt"
	"korero/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
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

var (
	table                 *tablewriter.Table
	rows                  chan []string
	channel               string
	currentSendingChannel string
)

func messagesFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	chanelFlagName := "channel"
	flags.StringVar(&channel, chanelFlagName, "", "only listen to messages in a certain channel (ID)")
}

func init() {
	messagesFlags(DiscordMessagesCommand)
	if len(channel) != 0 {
		currentSendingChannel = channel
	}
}

func messages(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		Token = args[0]
	}
	// login to the discord account given
	dg, err := utils.LoginDiscord(Token)
	if err != nil {
		return err
	}

	fmt.Println("to send a message to the channel in which messages are recieved, type and press enter")
	fmt.Printf("---------------------------------------------------------------------------------------\n\n")

	// create and set up channel writer for the message table
	rows = make(chan []string)
	table = tablewriter.NewWriter(os.Stdout)
	generateAscii(table)
	go func() {
		table.ContinuousRender(rows)
	}()

	// add discord event listeners
	dg.AddHandler(list)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		return err
	}

	foundChannel := false
	for _, guild := range dg.State.Guilds {
		channels, err := dg.GuildChannels(guild.ID)
		if err != nil {
			return err
		}
		for _, c := range channels {
			if c.Type != discordgo.ChannelTypeGuildText {
				continue
			}
			if c.Name == "general" {
				currentSendingChannel = c.ID
				foundChannel = true
				break
			}
			if foundChannel {
				break
			}
		}
	}

	go sendMessage(dg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()

	return nil
}

func list(dg *discordgo.Session, message *discordgo.MessageCreate) {
	t, err := message.Timestamp.Parse()
	if err != nil {
		os.Exit(125)
	}
	if message.Author.ID != dg.State.User.ID && (len(channel) == 0 || (message.ChannelID == channel)) {
		// write new content to the table
		rows <- []string{t.Format(time.RFC822), message.Content, message.Author.Username}
	}
}

// generateAscii created the headers for the korero message table
func generateAscii(table *tablewriter.Table) {
	table.SetHeader([]string{"Time", "Message", "User"})
}

// sendMessage is a goroutine which contantly is reading from stdin to see if the user is trying to send a message
func sendMessage(s *discordgo.Session) {
	for {
		reader := bufio.NewReader(os.Stdin)
		txt, _ := reader.ReadString('\n')
		if len(txt) > 0 {
			s.ChannelMessageSend(currentSendingChannel, txt)
		}
	}
}
