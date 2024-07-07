package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"module-go/internal/bot/commands/information"
	"module-go/internal/bot/commands/utilities"
	"module-go/internal/bot/handlers"
	"module-go/internal/bot/handlers/command"
	"module-go/internal/cfg"
	"module-go/internal/services"
	"module-go/internal/types"
)

func Start(guildService services.GuildService) {
	session, err := discordgo.New("Bot " + cfg.Get().DiscordToken)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	session.Identify.Intents = discordgo.IntentsAll
	session.StateEnabled = true
	session.State.MaxMessageCount = 100

	guildEvents := handlers.NewGuildEvents(guildService)

	commandHandler := command.NewCommandHandler(InitCommands(), guildService, cfg.Get().DiscordOwnerID)

	session.AddHandler(guildEvents.OnGuildCreate)
	session.AddHandler(commandHandler.OnInteractionCreate)

	if err := session.Open(); err != nil {
		log.Fatal().Err(err).Send()
	}

	RegisterCommands(session, commandHandler, cfg.Get().DiscordGuildID)

	log.Info().Msg("Bot started")

	select {}
}

func InitCommands() map[string]*command.Command {
	return map[string]*command.Command{
		"server": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Name:        "server",
				Description: "Information about server",
			},
			Category:          types.INFORMATION,
			OwnerCommand:      false,
			ModerationCommand: false,
			Hidden:            false,
			Handler:           &information.ServerCommand{},
		},
		"user": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Name:        "user",
				Description: "Information about user",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "Specific user",
					},
				},
			},
			Category:          types.INFORMATION,
			OwnerCommand:      false,
			ModerationCommand: false,
			Hidden:            false,
			Handler:           &information.UserCommand{},
		},
		"avatar": {
			ApplicationCommand: &discordgo.ApplicationCommand{
				Name:        "avatar",
				Description: "User avatar",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "Specific user",
					},
				},
			},
			Category:          types.UTILITIES,
			OwnerCommand:      false,
			ModerationCommand: false,
			Hidden:            false,
			Handler:           &utilities.AvatarCommand{},
		},
	}
}

func RegisterCommands(session *discordgo.Session, commandHandler *command.Handler, guildId string) {
	for _, cmd := range commandHandler.Commands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, guildId, cmd.ApplicationCommand)
		if err != nil {
			log.Error().Err(err).Send()
		}
	}
}
