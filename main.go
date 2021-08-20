package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	var err error
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatalf("Missing Discord authentication token. Check README on how to resolve this issue.")
	}
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error authenticating with Discord's servers. More information to follow: %v", err)
	}

	// Open connection to Discord
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot connect to Discord's servers. More information to follow: %v", err)
	}
	// Log OK and set status
	log.Println("=== === ===")
	log.Println("Bot is currently running.")
	log.Println("=== === ===")
	s.UpdateGameStatus(0, "Use v.help")

	s.AddHandler(cmd)
	s.AddHandler(reactAdd)

	// Gracefully close the Discord session, where possible
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stop
	s.Close()
	log.Println("Shutting down bot gracefully...")
}

func cmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Handling for each command there is
	if strings.HasPrefix(m.Content, "v.votemute") {
		if len(m.Mentions) >= 2 {
			// Notify that voting to mute several people isn't here.. yet.
			s.ChannelMessageSend(m.ChannelID, "Voting to mute several people isn't here.. yet. Check https://github.com/doamatto/vote-to-mute to see when this is added.")
		} else if len(m.Mentions) == 1 {
			// Mute only one user
			str := "Will we be muting " + m.Mentions[0].Mention() + " ? Vote on it!"
			msg, err := s.ChannelMessageSend(m.ChannelID, str)
			if err != nil {
				log.Panicf("%v", err)
			}

			// Add reaction
			err = s.MessageReactionAdd(m.ChannelID, msg.ID, "👍")
			if err != nil {
				log.Panicf("%v", err)
			}
		} else {
			// Notify that you must mention who to mute.
			s.ChannelMessageSend(m.ChannelID, "Please mention who to be muted (tip: type an @ followed by their name, or shift-click the user)")
		}
	}
	if strings.HasPrefix(m.Content, "v.about") {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       "About this bot",
			Color:       16724804,
			Description: "This was a bot written by [doamatto](https://www.doamatto.xyz) to both experiment with discordgo and help a friend with a moderation issue in a server.",
		})
	}
	if strings.HasPrefix(m.Content, "v.h") || strings.HasPrefix(m.Content, "v.help") {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: "Commands",
			Color: 16724804,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "v.about", Value: "What does this bot do and other FAQs", Inline: false},
				{Name: "v.votemute", Value: "Vote to mute whatever user you mention. Can't be someone with higher privileges than this bot.", Inline: false},
			},
		})
	}
}

func reactAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	m := r.MessageID
	c := r.ChannelID
	msg, err := s.ChannelMessage(c, m)
	if err != nil {
		log.Panicf("%v", err)
	}
	g := msg.GuildID

	// TODO: add a vote age limit (a day or so)

	// Ignore if message being reacted on isn't one of ours
	if msg.Author.ID != s.State.User.ID {
		return
	}

	// Ignore if emoji is not helping pass a vote
	if r.Emoji.Name != "👍" {
		return
	}

	// See if threshold is met
	if msg.Reactions[0].Emoji.Name == "👍" && msg.Reactions[0].Count >= 8 {
		// Fetch ID
		id := msg.Mentions[0].ID

		// Mute the user, if the role already exists
		roles, err := s.GuildRoles(g)
		if err != nil {
			log.Panicf("%v", err)
		}
		for _, r := range roles {
			if r.Name == "Muted" {
				// Give the user the Muted role
				s.GuildMemberRoleAdd(g, id, r.ID)
				return
			}
		}

		// Create the missing roles and give the role
		//
		// The role is grey; the user only gets permissions to read channels.
		// Servers will have to revoke permissions manually due to Discord not
		// giving privileges to allow these kinds of interactions (afaik).
		role, err := s.GuildRoleCreate(g)
		if err != nil {
			log.Panicf("%v", err)
		}
		s.GuildRoleEdit(g, role.ID, "Muted", 6052956, false, 66560, false)
	}
}
