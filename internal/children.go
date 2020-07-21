package internal

import (
	dg "github.com/bwmarrin/discordgo"
	"log"
)

func (bot *Bot) GetChildren(fluxCat FluxCategory) (children []*dg.Channel) {
	guild, err := bot.Client.State.Guild(bot.Serving)

	if err != nil {
		log.Println("Failed to channels of serving " + bot.Serving)
		return children
	}

	for _, channel := range guild.Channels {
		if channel.ParentID == fluxCat.CategoryID && channel.Type == dg.ChannelTypeGuildVoice {
			isAParent := false
			for _, parentID := range fluxCat.Parents {
				if parentID == channel.ID {
					isAParent = true
					break
				}
			}
			if !isAParent {
				children = append(children, channel)
			}
		}
	}
	return children
}

func (bot *Bot) GetChildrenCapacity(category FluxCategory) (childrenCapacity []Capacity) {
	children := bot.GetChildren(category)

	for _, child := range children {
		members := bot.GetMembersOfVC(child.ID)
		capacity := bot.BuildCapacity(members, child, false)
		childrenCapacity = append(childrenCapacity, capacity)
	}

	return childrenCapacity
}

func (bot *Bot) AreAllChildrenFull(category FluxCategory) (bool, []*dg.Channel) {
	childrenCapacity := bot.GetChildrenCapacity(category)
	var children []*dg.Channel
	var areAllFull = true

	for _, capacity := range childrenCapacity {
		children = append(children, capacity.Channel)

		if !capacity.IsFull {
			areAllFull = false
		}
	}

	return areAllFull, children
}

func (bot *Bot) GetEmptyChildren(category FluxCategory) (empty []*dg.Channel) {
	children := bot.GetChildrenCapacity(category)

	for _, child := range children {
		if child.IsEmpty {
			empty = append(empty, child.Channel)
		}
	}

	return empty
}

func (bot *Bot) DeleteEmptyChildren(category FluxCategory, areParentsFull bool) {
	var err error
	emptyChildren := bot.GetEmptyChildren(category)
	lastIndex := len(emptyChildren) - 1

	for i, child := range emptyChildren {
		// leave at least one empty child if the parents are full
		if i == lastIndex && !areParentsFull {
			_, err = bot.Client.ChannelDelete(child.ID)
			if err != nil {
				log.Printf("Failed to delete child channel %s, because\n%s", child.ID, err.Error())
			} else {
				log.Printf("Deleted %s's child %s\n", category.CategoryID, child.ID)
			}
		} else if i != lastIndex && areParentsFull {
			_, err = bot.Client.ChannelDelete(child.ID)
			if err != nil {
				log.Printf("Failed to delete child channel %s, because\n%s", child.ID, err.Error())
			} else {
				log.Printf("Deleted %s's child %s\n", category.CategoryID, child.ID)
			}
		}
	}
}
