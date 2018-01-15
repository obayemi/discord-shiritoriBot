package main

import (
    "flag"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/jinzhu/gorm"
    "github.com/bwmarrin/discordgo"
)


var (
    Db *gorm.DB
    DiscordToken string
    DiscordSession, _ = discordgo.New()
)


func initDiscord() {
    var err error
    if DiscordSession.Token == "" {
        fmt.Println("You must provide a Discord authentication token.")
        return
    }
    DiscordSession.State.User, err = DiscordSession.User("@me")
    if err != nil {
        fmt.Printf("error fetching user information, %s\n", err)
    }
    DiscordSession.AddHandler(MessageHandler)
    if err := DiscordSession.Open(); err != nil {
        fmt.Printf("error opening connection to Discord, %s\n", err)
        os.Exit(1)
    }
}


func init() {
    var discordToken string
    discordToken = os.Getenv("DISCORD_TOKEN")
    if DiscordSession.Token == "" {
        flag.StringVar(&discordToken, "t", "", "Discord Authentication Token")
    }
    flag.Parse()

    DiscordSession.Token = "Bot " + discordToken
}


func main() {
    Db = initDb()
    defer Db.Close()
    initDiscord()
    defer DiscordSession.Close()

    fmt.Println(`Now running. Press CTRL-C to exit.`)
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    var words []Word
    Db.Find(&words)
    fmt.Println("saved words")
    for i, word := range words {
        fmt.Printf("%d - %s\n", i, word.Word)
    }
    //Db.Delete(&words)
}
