package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unicode"
)

// needs 26688 perms

func ManipTxt(content string) string {
	var out []rune
	for i, c := range content {
		if i%2 == 0 {
			out = append(out, unicode.ToUpper(c))
		} else {
			out = append(out, unicode.ToLower(c))
		}
	}
	return string(out)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if seq := strings.Split(m.Content, "&sponge"); len(seq) > 1 {
		s.ChannelMessageSend(m.ChannelID, ManipTxt(seq[1]))
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func main() {
	viper.SetConfigName("discord")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	discord, err := discordgo.New("Bot " + viper.GetString("token"))
	discord.AddHandler(messageCreate)
	discord.Open()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
}
