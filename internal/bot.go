package internal

import (
	auth "github.com/Floor-Gang/authclient"
	util "github.com/Floor-Gang/utilpkg"
	dg "github.com/bwmarrin/discordgo"
	"log"
)

const (
	maxParentCapacity = 1
	maxChildCapacity  = 1
)

// Bot structure
type Bot struct {
	Auth          *auth.AuthClient
	Client        *dg.Session
	Config        *Config
	Serving       string // The guild we're serving
	OldVoiceState map[string]*dg.VoiceState
}

type FluxCategory struct {
	ChannelNamePrefix string   `yaml:"channel_name_prefix"`
	CategoryID        string   `yaml:"category_id"`
	Parents           []string `yaml:"parents"`
}

type Capacity struct {
	Members []string
	IsFull  bool
	IsEmpty bool
	Channel *dg.Channel
}

// Start starts discord client, configuration and database
func Start() {
	var err error

	// get Config.yml
	config := GetConfig()

	// setup authentication server
	// you can use this to get the bot's access token
	// and authenticate each user using a command.
	authClient, err := auth.GetClient(config.Auth)

	if err != nil {
		log.Fatalln("Failed to connect to authentication server", err)
	}

	_, err = authClient.Register(
		auth.Feature{
			Name:        "Flux Channels",
			Description: "Resize the number of voice channels in a category based on capacity",
			Commands: []auth.SubCommand{
				{
					Name:        "add",
					Description: "Add another category to resize automatically",
					Example:     []string{"add", "category ID"},
				},
				{
					Name:        "remove",
					Description: "Remove a fluctuating category",
					Example:     []string{"remove", "category ID"},
				},
				{
					Name:        "list",
					Description: "List all the fluctuating categories",
					Example:     []string{"list"},
				},
			},
			CommandPrefix: config.Prefix,
		},
	)

	if err != nil {
		log.Fatalln("Failed to register with authentication server", err)
	}

	client, err := dg.New("Bot " + config.Token)

	if err != nil {
		panic(err)
	}

	bot := Bot{
		Auth:          &authClient,
		Client:        client,
		Config:        &config,
		Serving:       config.Guild,
		OldVoiceState: make(map[string]*dg.VoiceState),
	}

	client.AddHandler(bot.onVoiceUpdate)
	client.AddHandler(bot.onMessage)
	client.AddHandlerOnce(bot.onReady)

	if err = client.Open(); err != nil {
		util.Report("Was an authentication token provided?", err)
	}
}

func (bot *Bot) GetMembersOfVC(voiceID string) (members []string) {
	guild, err := bot.Client.State.Guild(bot.Serving)

	if err != nil {
		log.Println("Failed to get serving guild " + bot.Serving)
		return
	}

	for _, voiceState := range guild.VoiceStates {
		if voiceState.ChannelID == voiceID {
			members = append(members, voiceState.UserID)
		}
	}

	return members
}
