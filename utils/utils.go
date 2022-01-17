package utils

import (
	"errors"
	"os"

	"github.com/bwmarrin/discordgo"
)

func LoginDiscord(token string) (*discordgo.Session, error) {
	if len(token) > 0 {
		if err := os.Setenv("KORERO_DISCORD_TOKEN", token); err != nil {
			return nil, err
		}
	} else {
		token = os.Getenv("KORERO_DISCORD_TOKEN")
	}
	if len(token) == 0 {
		return nil, errors.New("token length is zero, to use a default token please set $KORERO_DISCORD_TOKEN")
	}
	return discordgo.New("Bot " + token)
}
