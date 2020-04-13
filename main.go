package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unicode"
)

const Version = "0.4.1"

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

func ToClap(content string) string {
	var out []string
	for i, c := range strings.Split(content, " ") {
		if i == 0 {
			out = append(out, fmt.Sprintf("👏%s", c))
		} else {
			out = append(out, fmt.Sprintf("%s 👏", c))
		}
	}
	return strings.Join(out, " ")
}

func DiscordSetup() {
	// needs 26688 perms
	const TARGET_PERMS uint = 26688
	if x := viper.GetString("client_id"); x != "" {
		log.Printf("connect via https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=%d", x, TARGET_PERMS)
	} else {
		fmt.Println("What is the application's client_id (grab it from https://discordapp.com/developers/applications/ )?")
		var cl_id string
		fmt.Scanln(&cl_id)
		log.Printf("connect via https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=%d", cl_id, TARGET_PERMS)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if seq := strings.Split(m.Content, fmt.Sprintf("%ssponge ", viper.GetString("prefix"))); len(seq) > 1 {
		s.ChannelMessageSend(m.ChannelID, string(m.Author.Username)+"» "+ManipTxt(seq[1]))
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
	if seq := strings.Split(m.Content, fmt.Sprintf("%sclap ", viper.GetString("prefix"))); len(seq) > 1 {
		s.ChannelMessageSend(m.ChannelID, string(m.Author.Username)+"» "+ToClap(seq[1]))
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func main() {
	log.Printf("Version %v spongebob_exe bot", Version)
	log.Println("source at: https://github.com/jrabbit/spongebob_exe")
	viper.SetDefault("prefix", "&")
	viper.SetConfigName("spongebob")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("$HOME/.config/")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error in config file: %s \nDo you have a config?", err))
	}
	flag.Bool("setup", false, "give discord join url")
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)
	if viper.GetBool("setup") {
		DiscordSetup()
		return
	}
	log.Printf("using prefix %s", viper.GetString("prefix"))
	discord, err := discordgo.New("Bot " + viper.GetString("token"))
	discord.AddHandler(messageCreate)
	discord.Open()
	discord.UpdateListeningStatus(fmt.Sprintf("%ssponge <yr message>", viper.GetString("prefix")))
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
}
