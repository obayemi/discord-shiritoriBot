package main

import (
    "strings"
    "fmt"
    "github.com/bwmarrin/discordgo"
)


func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        return
    }
    commands := strings.Split(m.Content, " ")
    if commands[0] == "!shiritori" {
        interpretCommand(s, m, commands[1:])
        return
    }
    playShiritori(s, m)
}


func interpretCommand(
    s *discordgo.Session,
    m *discordgo.MessageCreate,
    commands []string,
) {
    switch commands[0] {
    case "start":
        initShiritori(s, m)
    case "init":
        initShiritori(s, m)
    case "backlog":
        backlogShiritori(s, m)
    case "stop":
        stopShiritori(s, m)
    case "reset":
        resetShiritori(s, m)
    case "clean":
        cleanShiritori(s, m)
    case "help":
        help(s, m)
    default:
        s.ChannelMessageSend(
            m.ChannelID,
            fmt.Sprintf("\"%s\" is not a valid shiritori command", commands[0]),
        )
    }
}
