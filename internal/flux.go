package internal

import (
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"log"
)

// This is how fluctuating channels works.
// first we need to describe things
// parent channels: Channels that will always exist.
// children channels: Channels that were created because the parents and the other
//                    existing children are full

// 1. The bot will be given a category and intake all the current existing voice channels and mark
// them as "parent channels". When these channels are full a new child channel is created
// 2. The bot will continue to monitor the capacity of all the parent and children channels.
// 3. If multiple children channels are empty then only one is kept and the rest are deleted

func (bot *Bot) CheckFluxCategory(category FluxCategory) {
	log.Printf("Checking fluctuating category %s\n", category.CategoryID)
	areParentsFull, parents := bot.AreAllParentsFull(category)
	areChildrenFull, children := bot.AreAllChildrenFull(category)

	if areParentsFull && areChildrenFull {
		bot.NewChild(category, parents, children)
		return
	}

	if !areParentsFull || !areChildrenFull {
		bot.DeleteEmptyChildren(category, areParentsFull)
		return
	}

	log.Println(category.CategoryID + " is all good")
}

func (bot *Bot) NewChild(category FluxCategory, parents []*dg.Channel, children []*dg.Channel) {
	log.Printf("Creating a new child for %s\n", category.CategoryID)
	var (
		newChild = dg.GuildChannelCreateData{
			Type:                 dg.ChannelTypeGuildVoice,
			Bitrate:              0,
			UserLimit:            0,
			PermissionOverwrites: nil,
			ParentID:             category.CategoryID,
		}
		size     = len(children)
		fullSize = size + len(parents)
	)

	if size > 0 {
		newChild.Name = fmt.Sprintf("%s %d", category.ChannelNamePrefix, fullSize+1)
		newChild.Position = children[size-1].Position + 1
	} else {
		newChild.Name = fmt.Sprintf("%s %d", category.ChannelNamePrefix, len(category.Parents)+1)
	}

	channel, err := bot.Client.GuildChannelCreateComplex(
		bot.Serving,
		newChild,
	)

	if err != nil {
		log.Printf(
			"Failed to create new child for %s because \n%s",
			category.CategoryID,
			err.Error(),
		)
	} else {
		log.Printf(
			"Created a new child for %s, it's ID is %s\n",
			category.CategoryID,
			channel.ID,
		)
	}
}

func (bot *Bot) BuildCapacity(members []string, channel *dg.Channel, isParent bool) (capacity Capacity) {
	size := len(members)
	capacity = Capacity{
		Members: members,
		IsFull:  false,
		IsEmpty: size == 0,
		Channel: channel,
	}

	if isParent {
		capacity.IsFull = size >= maxParentCapacity
	} else {
		capacity.IsFull = size >= maxChildCapacity
	}

	return capacity
}
