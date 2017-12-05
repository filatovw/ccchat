package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

func GetOrCreateUser(db *gorm.DB, name string) (*User, error) {
	user := &User{Name: name}
	res := db.FirstOrCreate(user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func AddMessage(db *gorm.DB, user *User, body string) error {
	m := &Message{Body: body, UserID: user.ID}
	res := db.Model(user).Association("Messages").Append(m)
	return res.Error
}

func MessagesSinceDate(db *gorm.DB, t time.Time) ([]Message, error) {
	msgs := []Message{}
	res := db.Where("created_at > ?", t).Find(&msgs)
	if res.Error != nil {
		return nil, res.Error
	}
	return msgs, nil
}
