package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"../config"
	"../handlers"
)

//BotID is bot id
var BotID string
var goBot *discordgo.Session

//Start starting the bot
func Start() {

	fmt.Println(config.BotName + " is starting...")
	defer fmt.Println(config.BotName + " is running")

	var err error
	goBot, err = discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

//Stop stopping the bot
func Stop() {
	goBot.Close()
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotID {
		return
	}

	go func() {
		if strings.HasPrefix(m.Content, config.BotPrefix) {

			if m.Content == "!ping" {
				_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
			}

			//roll dices
			if strings.HasPrefix(m.Content, "!r") || strings.HasPrefix(m.Content, "!roll") {

				fmt.Println("command:", m.Content)
				var answer string

				//roll with verbose
				if strings.HasPrefix(m.Content, "!rv") || strings.HasPrefix(m.Content, "!rollv") {
					content := strings.TrimPrefix(m.Content, "!rollv ")
					content = strings.TrimPrefix(content, "!rv ")

					answer = handlers.Roll(content, true)
				} else { //without verbose
					content := strings.TrimPrefix(m.Content, "!roll ")
					content = strings.TrimPrefix(content, "!r ")

					answer = handlers.Roll(content, false)
				}

				answer = "```" + answer + "```"

				_, _ = s.ChannelMessageSend(m.ChannelID, answer)
			}

			//check online minecraft server
			if strings.HasPrefix(m.Content, "!online") {

				fmt.Println("command:", m.Content)

				content := strings.TrimPrefix(m.Content, "!online")

				answer := handlers.PingMinecraftServer(content)

				_, _ = s.ChannelMessageSend(m.ChannelID, answer)

				fmt.Println(answer)
			}

		}
	}()

}
