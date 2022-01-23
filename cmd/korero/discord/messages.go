package discord

import (
	"korero/utils"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var (
	messagesDescription = `logs in to the default user and streams messages from their servers`

	DiscordMessagesCommand = &cobra.Command{
		Use:     "messages TOKEN",
		Short:   "stream messages from the given user",
		Long:    messagesDescription,
		RunE:    messages,
		Example: "korero discord messages 12345",
	}
)

var (
	rows                  chan []string
	allRows               []string
	channel               string
	user                  string
	originalUser          string
	currentSendingChannel string
	app                   *tview.Application
	table                 *tview.Table
	grid                  *tview.Grid
	form                  *tview.Form
	tree                  *tview.TreeView
	serverList            []*discordgo.Guild
	channelList           []*discordgo.Channel
)

func messagesFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	userFlagName := "user"
	flags.StringVar(&user, userFlagName, "", "set the discord bot to use a specific name when sending messages")

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
	// create writer for the message table
	rows = make(chan []string)

	// add discord event listeners
	dg.AddHandler(list)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		return err
	}

	if len(user) > 0 {
		originalUser = dg.State.User.Username
		_, err := dg.UserUpdate("", "", user, "", "")
		if err != nil {
			os.Exit(125)
		}
	}

	foundChannel := false
	for _, guild := range dg.State.Guilds {
		serverList = append(serverList, guild)
		channels, err := dg.GuildChannels(guild.ID)
		if err != nil {
			return err
		}
		for _, c := range channels {
			channelList = append(channelList, c)
			if c.Type != discordgo.ChannelTypeGuildText {
				continue
			}
			if c.Name == "general" { // just join a channel by default
				currentSendingChannel = c.ID
				foundChannel = true
				break
			}
			if foundChannel {
				break
			}
		}
	}

	generateView(dg)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			if len(originalUser) > 0 {
				_, err := dg.UserUpdate("", "", originalUser, "", "")
				if err != nil {
					os.Exit(125)
				}
			}
			app.Stop()
			os.Exit(0)

		case tcell.KeyEnter:
			sendMessage(dg, form.GetFormItem(0).(*tview.InputField).GetText())
			form.GetFormItem(0).(*tview.InputField).SetText("")
		}
		return event
	})

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		if node.GetText() != "." {
			val := reflect.ValueOf(node.GetReference()).Elem()
			id := val.FieldByName("ID").Interface().(string)
			currentSendingChannel = id
		}
	})

	go func() {
		updateTable(rows)
	}()
	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()

	return nil
}

// list is the function to add new messages to the rows of the table
func list(dg *discordgo.Session, message *discordgo.MessageCreate) {
	t, err := message.Timestamp.Parse()
	t = t.In(time.Local)
	if err != nil {
		os.Exit(125)
	}
	if message.Author.ID != dg.State.User.ID && (len(currentSendingChannel) == 0 || (message.ChannelID == currentSendingChannel)) {
		// write new content to the table
		rows <- []string{t.Format(time.Kitchen), message.Content, message.Author.Username}

	}
}

// generateView creates the entire grid view
func generateView(dg *discordgo.Session) {
	app = tview.NewApplication()
	table = tview.NewTable().
		SetBorders(true) // init table
	grid = tview.NewGrid().SetRows(5, 0, 3).SetColumns(30, 0, 30).SetBorders(true)
	form = tview.NewForm().AddInputField("Send:", "", 30, nil, nil).AddButton("Send Message", func() {
		sendMessage(dg, form.GetFormItem(0).(*tview.InputField).GetText())
		form.GetFormItem(0).(*tview.InputField).SetText("")
	}) // create form and add function for when the button is clicked
	grid.AddItem(tview.NewTextView().SetText( // title header
		"    __ __\n"+
			"   / //_/___  ________ _________ \n"+
			"  / ,< / __ \\/___/ _ \\/ ___/ __ \\\n"+
			" / /| / /_/ / / /  __/ /  / /_/ /\n"+
			"/_/ |_\\____/_/  \\___/_/   \\____/\n\n"), 0, 0, 1, 3, 0, 0, false)
	root := tview.NewTreeNode(".").
		SetColor(tcell.ColorRed)
	tree = tview.NewTreeView().SetRoot(root).
		SetCurrentNode(root)
	grid.AddItem(form, 1, 0, 1, 1, 0, 50, false)                                                                                                                              // add the form to the grid
	grid.AddItem(table, 1, 1, 1, 1, 0, 50, false)                                                                                                                             // add the table to the grid
	grid.AddItem(tview.NewTextView().SetText("<ENTER> Send Message \t <CTR Q> Exit \t <ARROW KEYS> Navigate lists").SetTextAlign(tview.AlignCenter), 2, 0, 1, 3, 0, 0, false) // add command bar to the grid

	for _, server := range serverList {
		s, err := dg.Guild(server.ID)
		if err != nil {
			os.Exit(125)
		}
		nodeServer := tview.NewTreeNode(s.Name).
			SetReference(s).SetSelectable(false)
		root.AddChild(nodeServer).SetColor(tcell.ColorGreen)
		channels, err := dg.GuildChannels(server.ID)
		if err != nil {
			os.Exit(125)
		}
		for _, channel := range channels {
			if channel.Type != discordgo.ChannelTypeGuildCategory && channel.Type != discordgo.ChannelTypeGuildVoice {
				node := tview.NewTreeNode(channel.Name).
					SetReference(channel).SetSelectable(true)
				nodeServer.AddChild(node)
			}
		}

	}
	grid.AddItem(tree, 1, 2, 1, 1, 0, 50, false).SetTitle("Server List") // add server list to the grid

}

// sendMessage reads from the inputField and sends the message to the discord channel selected
func sendMessage(sess *discordgo.Session, txt string) {
	if len(txt) > 0 {
		sess.ChannelMessageSend(currentSendingChannel, txt)
	}
}

// updateTable is called whenever the row channel is given a new entry, it redraws the entire table with the new entries.
func updateTable(rows <-chan []string) {
	for {
		word := 0
		updatedRows := <-rows
		allRows = append(allRows, updatedRows...)
		app.QueueUpdateDraw(func() {
			table.SetCell(0, 0, tview.NewTableCell("TIME").
				SetTextColor(tcell.ColorBlue).
				SetAlign(tview.AlignCenter))
			table.SetCell(0, 1, tview.NewTableCell("MESSAGE").
				SetTextColor(tcell.ColorBlue).
				SetAlign(tview.AlignCenter))
			table.SetCell(0, 2, tview.NewTableCell("USER").
				SetTextColor(tcell.ColorBlue).
				SetAlign(tview.AlignCenter))
			for r := 1; r <= ((len(allRows) + 1) / 3); r++ {
				for c := 0; c < 3; c++ {
					color := tcell.ColorWhite
					if c < 1 || r < 1 {
						color = tcell.ColorYellow
					}
					table.SetCell(r, c,
						tview.NewTableCell(allRows[word]).
							SetTextColor(color).
							SetAlign(tview.AlignCenter).SetSelectable(true))
					word = (word + 1) % len(allRows)
				}
			}
			table.ScrollToEnd()
		})
	}
}
