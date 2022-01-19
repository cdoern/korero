package discord

import (
	"errors"
	"fmt"
	"korero/utils"
	"os"

	"github.com/spf13/cobra"
)

var (
	loginDescription = `Creates a new login instance for the discord API`

	DiscordLoginCommand = &cobra.Command{
		Use:     "login TOKEN",
		Short:   "login to the given user",
		Long:    loginDescription,
		RunE:    login,
		Example: "korero discord login 12345",
	}
)

func loginFlags() {

}

func init() {
	loginFlags()
}

func login(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		Token = args[0]
		if err := os.Setenv("KORERO_DISCORD_TOKEN", Token); err != nil {
			return err
		}
	} else {
		Token = os.Getenv("KORERO_DISCORD_TOKEN")
	}

	if len(Token) == 0 {
		return errors.New("invalid Argument, token length is zero")
	}
	dg, err := utils.LoginDiscord(Token)
	if err != nil {
		return err
	}
	fmt.Println("Login Successful:", dg.Identify.Token)
	return nil
}
