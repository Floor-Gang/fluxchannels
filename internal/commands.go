package internal

import (
	"fmt"
	util "github.com/Floor-Gang/utilpkg/botutil"
	dg "github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// args = [prefix, add, category ID, channel name prefix...]
func (bot *Bot) cmdAdd(msg *dg.Message, args []string) {
	if len(args) < 4 {
		_, _ = util.Reply(bot.Client, msg, bot.Config.Prefix+" add <category ID> <channel name prefix>")
		return
	}

	categoryID := args[2]
	category, err := bot.Client.Channel(categoryID)
	channelNamePrefix := strings.Join(args[3:], " ")
	log.Println(channelNamePrefix)

	if err != nil || category.Type != dg.ChannelTypeGuildCategory {
		_, _ = util.Reply(bot.Client, msg, fmt.Sprintf("`%s` isn't an ID of a category.", categoryID))
		return
	}

	guild, _ := bot.Client.State.Guild(bot.Serving)

	if guild == nil {
		return
	}
	var parents []string
	var result = fmt.Sprintf("**Added Category** %s\nwith parents\n", category.Name)
	for _, channel := range guild.Channels {
		if channel.Type == dg.ChannelTypeGuildVoice && channel.ParentID == category.ID {
			result += fmt.Sprintf(" - %s\n", channel.Name)
			parents = append(parents, channel.ID)
		}
	}

	fluxCategory := FluxCategory{
		CategoryID:        category.ID,
		Parents:           parents,
		ChannelNamePrefix: channelNamePrefix,
	}
	bot.Config.Categories[category.ID] = fluxCategory
	if err := bot.Config.Save(); err != nil {
		_, _ = util.Reply(bot.Client, msg, "Failed to add category, something went wrong.")
		log.Println(err)
	} else {
		_, _ = util.Reply(bot.Client, msg, result)
	}
}

func (bot *Bot) cmdRemove(msg *dg.Message, args []string) {
	if len(args) < 3 {
		_, _ = util.Reply(bot.Client, msg, bot.Config.Prefix+" remove <category ID>")
		return
	}

	categoryID := args[2]
	category, isOK := bot.Config.Categories[categoryID]

	if !isOK {
		_, _ = util.Reply(bot.Client, msg, fmt.Sprintf("`%s` isn't an ID of a category.", categoryID))
		return
	}

	delete(bot.Config.Categories, category.CategoryID)
	if err := bot.Config.Save(); err != nil {
		_, _ = util.Reply(bot.Client, msg, "Failed to remove category, something went wrong.")
		log.Println(err)
	} else {
		_, _ = util.Reply(bot.Client, msg, "Removed category.")
	}
}

func (bot *Bot) cmdList(msg *dg.Message) {
	var list = "**Fluctuating Channels**\n"

	if len(bot.Config.Categories) > 0 {
		list += "```\n"
		for _, fluxCat := range bot.Config.Categories {
			category, err := bot.Client.Channel(fluxCat.CategoryID)
			if err != nil {
				list += fmt.Sprintf("Category - unknown (`%s`)", fluxCat.CategoryID)
			} else {
				list += "Category - " + category.Name
			}
			list += "\n"
			list += fmt.Sprintf("With channel prefix \"%s\"\n", fluxCat.ChannelNamePrefix)

			list += " Parent Channels\n"
			parents := bot.GetParents(fluxCat)
			if len(parents) > 0 {
				for _, parent := range parents {
					list += "  * " + parent.Name + "\n"
				}
			} else {
				list += "  None."
			}

			list += " Children Channels\n"
			children := bot.GetChildren(fluxCat)
			if len(children) > 0 {
				for _, child := range children {
					list += "  * " + child.Name + "\n"
				}
			} else {
				list += "    None"
			}
		}
		list += "```"
	} else {
		list += "There aren't any!"
	}

	_, _ = util.Reply(bot.Client, msg, list)
}
