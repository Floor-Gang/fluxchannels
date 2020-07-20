package internal

import (
	util "github.com/Floor-Gang/utilpkg/botutil"
	dg "github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func (bot *Bot) onMessage(_ *dg.Session, msg *dg.MessageCreate) {
	// check if they provided a prefix and they're not a bot
	if msg.Author.Bot || !strings.HasPrefix(msg.Content, bot.Config.Prefix) {
		return
	}

	// we can ask the authentication server if they're an admin of the bot
	// isAdmin defaults to "false" if err is not nil
	isAdmin, err := bot.Auth.CheckMember(msg.Author.ID)

	// let's split their message up into arguments
	// args = [prefix, sub-command name]
	args := strings.Fields(msg.Content)

	if len(args) < 2 { // this would mean args is [prefix] which at that point ignore them
		return
	}

	// we can now find out what command they were calling
	switch args[1] {
	case "ping":
		bot.cmdPing(msg.Message)
		break
	case "some-admin-command":
		// if all your commands are admin-relate then just wrap the whole switch statement
		// with this if statement

		// you can also just do "if isAdmin", the error doesn't matter that much if you want cleaner
		// code.
		if err != nil {
			util.Reply(bot.Client, msg.Message, "Failed to contact auth server")
		} else if isAdmin {
			bot.cmdSomethingAdmin(msg.Message)
		} else {
			util.Reply(bot.Client, msg.Message, "You must be an admin to run this command.")
		}
	}
}

func (bot *Bot) onReady(_ *dg.Session, ready *dg.Ready) {
	log.Printf("<Bot/Feature Name> - ready as %s#%s", ready.User.Username, ready.User.Discriminator)
}
