package main

import (
    "fmt"
    "os"
    "strings"
    "github.com/bwmarrin/discordgo"
)


func initShiritori(s *discordgo.Session, m *discordgo.MessageCreate) {
    var channel Channel
    Db.Where(Channel{
        ChannelID: m.ChannelID,
    }).Attrs(Channel{
        Active: true,
    }).FirstOrInit(&channel)

    if channel.Active && !Db.NewRecord(channel) {
        s.ChannelMessageSend(
            m.ChannelID,
            "shiritori session already in progress",
        )
        return
    }
    if Db.NewRecord(channel) {
        Db.Create(&channel)
    } else if channel.Active == false {
        channel.Active = true
        Db.Model(&channel).Update(Channel{Active: true})
    }
        s.ChannelMessageSend(
            m.ChannelID,
            fmt.Sprintf(
                "%s started a new shiritori session in <#%s>",
                m.Author.Mention(),
                m.ChannelID,
            ),
        )
}


func backlogShiritori(s *discordgo.Session, m *discordgo.MessageCreate) {
    sessions, err := s.ChannelMessages(m.ChannelID, 100, m.ID, "", "")
    if err != nil {
        s.ChannelMessageSend(
            m.ChannelID, "an error occured trying to fetch previous messages",
        )
    } else {
        for _, session := range sessions {
            if !(
                session.Author.ID == s.State.User.ID ||
                strings.HasPrefix(session.Content, "!")) {
                    Db.FirstOrCreate(&Word{
                        Word: session.Content,
                        ChannelID: session.ChannelID,
                        Author: session.Author.ID,
                    })
            }
        }
        s.ChannelMessageSend(
            m.ChannelID, "finised backlogging shiritori",
        )
    }
}


func stopShiritori(s *discordgo.Session, m *discordgo.MessageCreate) {
    var channel Channel
    err := Db.First(&channel, &Channel{ChannelID:m.ChannelID}).Error
    if err != nil || channel.Active == false {
        s.ChannelMessageSend(
            m.ChannelID,
            fmt.Sprintf(
                "no active shiritori in <#%s>",
                m.ChannelID))
    } else {
        tx := Db.Begin()
        tx.Delete(&Word{ChannelID: channel.ChannelID})
        tx.Delete(&channel)
        tx.Commit()
        s.ChannelMessageSend(
            m.ChannelID,
            fmt.Sprintf(
                "stopped shiritori session in <#%s>",
                m.ChannelID))
    }
}


func help(s *discordgo.Session, m *discordgo.MessageCreate) {
    s.ChannelMessageSend(
        m.ChannelID, "https://youtu.be/2dbR2JZmlWo",
    )
}


func resetShiritori(s *discordgo.Session, m *discordgo.MessageCreate) {
    var words []Word
    Db.Where(&Word{ChannelID: m.ChannelID}).Find(&words)

    for i, word := range words {
        fmt.Printf("%d - %s\n", i, word.Word)
    }
    Db.Delete(&words)

    s.ChannelMessageSend(
        m.ChannelID,
        "words DB emptied",
    )
    return
}


func cleanShiritori(s *discordgo.Session, m *discordgo.MessageCreate) {

    var channel Channel
    if err := Db.Where(&channel, ).Error; err != nil {
        s.ChannelMessageSend(
            m.ChannelID,
            "Can't clean a channel which is not a shiritori channel",
        )
        return
    } else {
        if false {
            sessions, err := s.ChannelMessages(m.ChannelID, 100, m.ID, "", "")
            if err != nil {
                s.ChannelMessageSend(
                    m.ChannelID, "an error occured trying to fetch previous messages",
                )
            } else {
                ids := make([]string, len(sessions))
                for i, s := range sessions {
                    ids[i] = s.ID
                }

                s.ChannelMessagesBulkDelete(m.ChannelID, ids)
            }
        } else if greatContent, err := os.Open("./fuck.gif"); err != nil {
            s.ChannelMessageSend(m.ChannelID, "not implemented")
        } else {
            s.ChannelFileSendWithMessage(m.ChannelID, "not implemented", "fuckyou.gif", greatContent)
        }
    }
}


func playShiritori(s *discordgo.Session, m *discordgo.MessageCreate) {
    var channel Channel
    if err := Db.Find(&channel, &Channel{ChannelID: m.ChannelID}).Error; err != nil {
        return
    }

    var word Word
    if err := Db.First(&word, &Word{Word: m.Content, ChannelID: channel.ChannelID}).Error; err == nil {
        user, _ := DiscordSession.User(word.Author)
        s.ChannelMessageSend(
            m.ChannelID,
            fmt.Sprintf("word \"%s\" already entered by %s", word.Word, user.Mention()),
        )
    } else {
        Db.Create(&Word{
            Word: m.Content,
            Author: m.Author.ID,
            ChannelID: channel.ChannelID,
        })
        fmt.Printf("[new word] %s\n", m.Content)
    }
}
