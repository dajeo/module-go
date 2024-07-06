package command

import (
	"github.com/bwmarrin/discordgo"
	"module-go/internal/types"
)

type Context struct {
	session *discordgo.Session
	event   *discordgo.InteractionCreate
	command *Command
}

func (ctx *Context) Reply(embed *discordgo.MessageEmbed) error {
	return ctx.session.InteractionRespond(ctx.event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func (ctx *Context) ReplyError(message string) error {
	embed := &discordgo.MessageEmbed{Description: message, Color: types.ERROR.Int()}
	return ctx.Reply(embed)
}

func (ctx *Context) Option(key string) *discordgo.ApplicationCommandInteractionDataOption {
	opts := ctx.event.ApplicationCommandData().Options
	for _, opt := range opts {
		if opt.Name == key {
			return opt
		}
	}
	return nil
}

func (ctx *Context) OptionAsUser(key string, defaultUser ...*discordgo.User) *discordgo.User {
	opt := ctx.Option(key)
	if opt != nil {
		return opt.UserValue(ctx.session)
	}

	if len(defaultUser) > 0 {
		return defaultUser[0]
	}

	return nil
}

func (ctx *Context) Guild() (*discordgo.Guild, error) {
	guild, err := ctx.session.State.Guild(ctx.event.GuildID)
	if err != nil {
		guild, err = ctx.session.Guild(ctx.event.GuildID)
	}
	return guild, err
}

func (ctx *Context) Member() *discordgo.Member {
	return ctx.event.Member
}

func (ctx *Context) User() *discordgo.User {
	return ctx.event.Member.User
}

func (ctx *Context) MemberByID(id string) (*discordgo.Member, error) {
	member, err := ctx.session.State.Member(ctx.event.GuildID, id)
	if err != nil {
		member, err = ctx.session.GuildMember(ctx.event.GuildID, id)
	}
	return member, err
}

func (ctx *Context) GuildOwner(ownerId string) (*discordgo.Member, error) {
	return ctx.MemberByID(ownerId)
}
