package discord

import (
	"bufio"
	"fmt"
	"korero/utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	setupDescription = `Helps user set up their discord configuration`

	DiscordSetupCommand = &cobra.Command{
		Use:     "setup [options]",
		Short:   "set up discord API connection",
		Long:    setupDescription,
		RunE:    setup,
		Example: "korero discord login 12345",
	}
)
var (
	Token string
)

func setupFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	tokenFlagName := "token"
	flags.StringVar(&Token, tokenFlagName, "", "set up the client with a token")
}

func init() {
	setupFlags(DiscordSetupCommand)
}

func setup(cmd *cobra.Command, args []string) error {
	if len(Token) == 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Welcome to Korero, the multipurpose lightweight messaging service.\n" +
			"are you setting up a discord API client? (y/n): ")
		response, _ := reader.ReadString('\n')
		evalResponse(response, "y\n")
		fmt.Print("\n1) go to https://discord.com/developers/applications and click \"New Application\" give your bot a name! type c to to continue. ")
		response, _ = reader.ReadString('\n')
		evalResponse(response, "c\n")
		fmt.Print("\n2) Navigate to the side bar and click the \"Bot\" tab, then click \"Add Bot\" and confirm. type c to continue. ")
		response, _ = reader.ReadString('\n')
		evalResponse(response, "c\n")
		fmt.Print("\n3) Make sure your bot is public, copy the bot's token and enter it here: ")
		Token, _ = reader.ReadString('\n')
		Token = strings.Replace(Token, "\n", "", -1)
		os.Setenv("KORERO_DISCORD_TOKEN", Token)
		fmt.Println("Using Provided Token...")
		dg, err := utils.LoginDiscord(Token)
		if err != nil {
			return err
		}
		fmt.Println("\nLogin Successful: ", dg.Identify.Token)
		fmt.Print("\n4) Great! now to invite your bot, navitage to the OAuth2 General tab on the webpage. Under \"Default Authorization Link\" select \"In-app Authorization\". type c to continue. ")
		response, _ = reader.ReadString('\n')
		evalResponse(response, "c\n")
		fmt.Print("\n5) Under Scopes click \"bot\" and \"applications.commands\". In the bot permissions table, select all of the commands you wish your bot to have! type c to continue. ")
		response, _ = reader.ReadString('\n')
		evalResponse(response, "c\n")
		fmt.Print("\n6) Navigate to the OAuth2 URL Generator. Again select \"bot\" and \"applications.commands\" then choose your permissions again to generate the URL. This url can then be used to add your bot to any server you own or can be distributed to any server who you know the owner of. type c to finish setup. ")
		response, _ = reader.ReadString('\n')
		evalResponse(response, "c\n")
	} else {
		fmt.Println("\nUsing provided token....")
		dg, err := utils.LoginDiscord(Token)
		if err != nil {
			return err
		}
		fmt.Println("\nLogin Successful: ", dg.Identify.Token)
	}
	return nil
}

func evalResponse(response string, lookingFor string) {
	switch response == lookingFor {
	case true:
		return
	case false:
		os.Exit(0)
	}

}
