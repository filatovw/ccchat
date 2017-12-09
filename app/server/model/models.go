package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	Name     string `gorm:"unique_index"`
	Messages []Message
}

type Message struct {
	gorm.Model
	Body   string
	UserID uint `gorm:"index"`
}

type UserMessage struct {
	UserName         string
	MessageCreatedAt time.Time
	MessageBody      string
}

func InitDB(host, user, password, database string) (*gorm.DB, error) {
	conn, err := gorm.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, database, password))
	if err != nil {
		return nil, err
	}

	conn.AutoMigrate(&User{}, &Message{})
	return conn, nil
}
