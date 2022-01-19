package discord

import (
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
	table *tablewriter.Table
	rows  chan []string
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
	// login to the discord account given
	dg, err := utils.LoginDiscord(Token)
	if err != nil {
		return err
	}

	// create and set up channel writer for the message table
	rows = make(chan []string)
	table = tablewriter.NewWriter(os.Stdout)
	generateAscii(table)
	go func() {
		table.ContinuousRender(rows)
	}()

	// add discord event listeners
	dg.AddHandler(list)

	//dg.AddHandler(wait)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		return err
	}
	input := make(chan *os.File)
	go getInput(input)
	go sendMessage(dg, input)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	//fmt.Println("here")

	dg.Close()

	return nil
}

func list(dg *discordgo.Session, message *discordgo.MessageCreate) {
	t, err := message.Timestamp.Parse()
	if err != nil {
		os.Exit(125)
	}
	// write new content to the table
	rows <- []string{t.Format(time.RFC822), message.Content, message.Author.Username}
}

func generateAscii(table *tablewriter.Table) {
	table.SetHeader([]string{"Time", "Message", "User"})
}

func getInput(input chan<- *os.File) {
	for {
		fmt.Println("here")
		time.Sleep(time.Second * 8)
		input <- os.Stdin
	}
}

func sendMessage(s *discordgo.Session, input chan *os.File) {
	for {
		msg := <-input
		s.ChannelMessageSend("369139060750417925", msg.Name())
	}
}
