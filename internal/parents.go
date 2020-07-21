package internal

import (
	dg "github.com/bwmarrin/discordgo"
	"log"
)

func (bot *Bot) AreAllParentsFull(category FluxCategory) (bool, []*dg.Channel) {
	var (
		parentCapacity = bot.GetParentCapacity(category)
		channels       []*dg.Channel
		areAllFull     = true
	)

	for _, capacity := range parentCapacity {
		channels = append(channels, capacity.Channel)

		if !capacity.IsFull {
			areAllFull = false
		}
	}

	return areAllFull, channels
}

func (bot *Bot) GetParents(category FluxCategory) (parents []*dg.Channel) {
	for _, channelID := range category.Parents {
		channel, err := bot.Client.Channel(channelID)

		if err != nil {
			log.Printf("Failed to get parent %s because,\n"+err.Error(), channelID)
		} else {
			parents = append(parents, channel)
		}
	}
	return parents
}

func (bot *Bot) GetParentCapacity(category FluxCategory) (parentCapacity []Capacity) {
	parents := bot.GetParents(category)

	for _, parent := range parents {
		members := bot.GetMembersOfVC(parent.ID)
		capacity := bot.BuildCapacity(members, parent, true)
		parentCapacity = append(parentCapacity, capacity)
	}
	return parentCapacity
}
