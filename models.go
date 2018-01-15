package main

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)


type Word struct {
    gorm.Model
    Word string
    Author string
    ChannelID string
}


type Channel struct {
    gorm.Model
    ChannelID string
    Active bool
}


func initDb() *gorm.DB {
    // func to init the DB, don't forget to close connection
    Db, err := gorm.Open("sqlite3", "shiritori.db")
    if err != nil {
        panic("failed to connect database")
    }
    // Migrate the schema
    Db.AutoMigrate(&Word{}, &Channel{})

    return Db
}
