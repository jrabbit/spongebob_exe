package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unicode"
)

// needs 26688 perms

const Version = "0.3.0"

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
	if seq := strings.Split(m.Content, fmt.Sprintf("%ssponge", viper.GetString("prefix"))); len(seq) > 1 {
		s.ChannelMessageSend(m.ChannelID, string(m.Author.Username)+ManipTxt(seq[1]))
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func main() {
	log.Printf("Version %v spongebob_exe bot", Version)
	log.Println("source at: https://github.com/jrabbit/spongebob_exe")
	viper.SetDefault("prefix", "&")
	viper.SetConfigName("discord")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error in config file: %s \nDo you have a config?", err))
	}
	discord, err := discordgo.New("Bot " + viper.GetString("token"))
	discord.AddHandler(messageCreate)
	discord.Open()
	discord.UpdateListeningStatus(fmt.Sprintf("%ssponge <yr message>", viper.GetString("prefix")))
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
}
