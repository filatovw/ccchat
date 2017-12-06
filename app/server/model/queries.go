package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// GetOrCreateUser find user by name or create
func GetOrCreateUser(db *gorm.DB, name string) (*User, error) {
	user := &User{Name: name}
	res := db.FirstOrCreate(user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

// AddMessage can add message to any user
func AddMessage(db *gorm.DB, user *User, body string) error {
	m := &Message{Body: body, UserID: user.ID}
	res := db.Model(user).Association("Messages").Append(m)
	return res.Error
}

// MessagesSinceDate should return list of messages crafted after date
func MessagesSinceDate(db *gorm.DB, t time.Time) ([]Message, error) {
	msgs := []Message{}
	res := db.Where("created_at > ?", t).Find(&msgs)
	if res.Error != nil {
		return nil, res.Error
	}
	return msgs, nil
}