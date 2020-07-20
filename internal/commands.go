package internal

import (
	util "github.com/Floor-Gang/utilpkg/botutil"
	dg "github.com/bwmarrin/discordgo"
)

// Example command
func (bot *Bot) cmdPing(msg *dg.Message) {
	util.Reply(bot.Client, msg, "Pong!")
}

func (bot *Bot) cmdSomethingAdmin(msg *dg.Message) {
	util.Reply(bot.Client, msg, "Some admin command")
}
