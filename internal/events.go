package internal

import (
	util "github.com/Floor-Gang/utilpkg/botutil"
	dg "github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func (bot *Bot) onMessage(_ *dg.Session, msg *dg.MessageCreate) {
	if msg.Author.Bot || !strings.HasPrefix(msg.Content, bot.Config.Prefix) {
		return
	}

	if msg.GuildID != bot.Serving {
		return
	}

	isAdmin, _ := bot.Auth.CheckMember(msg.Author.ID)
	// args = [prefix, sub-command]
	args := strings.Fields(msg.Content)

	if len(args) < 2 {
		return
	}

	if isAdmin {
		switch args[1] {
		case "add":
			bot.cmdAdd(msg.Message, args)
			break
		case "remove":
			bot.cmdRemove(msg.Message, args)
			break
		case "list":
			bot.cmdList(msg.Message)
			break
		}
	} else {
		util.Reply(bot.Client, msg.Message, "You must be an administrator to run this command")
	}
}

func (bot *Bot) onReady(_ *dg.Session, ready *dg.Ready) {
	log.Printf("Flux Channels - ready as %s#%s", ready.User.Username, ready.User.Discriminator)
	for _, category := range bot.Config.Categories {
		bot.CheckFluxCategory(category)
	}
}

func (bot *Bot) onVoiceUpdate(_ *dg.Session, voice *dg.VoiceStateUpdate) {
	if oldState, isOK := bot.OldVoiceState[voice.UserID]; isOK {
		bot.checkState(oldState)
	}
	bot.OldVoiceState[voice.UserID] = voice.VoiceState

	bot.checkState(voice.VoiceState)
}

func (bot *Bot) checkState(state *dg.VoiceState) {
	channel, err := bot.Client.State.Channel(state.ChannelID)
	if err != nil {
		return
	}
	if category, isOK := bot.Config.Categories[channel.ParentID]; isOK {
		bot.CheckFluxCategory(category)
	}
}
