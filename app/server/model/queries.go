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
	res := db.Model(user).Association("Messages").Append(&Message{Body: body})
	return res.Error
}

// MessagesSinceDate should return list of messages crafted since presented date
func MessagesSinceDate(db *gorm.DB, t time.Time) ([]UserMessage, error) {
	msgs := []UserMessage{}
	rows, err := db.Table("users").Select(
		"users.name, messages.created_at, messages.body").Joins(
		"join messages on users.id = messages.user_id ").Where("messages.created_at >= ?", t).Order("messages.created_at asc").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var name string
	var created time.Time
	var body string
	for rows.Next() {
		rows.Scan(&name, &created, &body)
		msgs = append(msgs, UserMessage{UserName: name, MessageCreatedAt: created, MessageBody: body})
	}
	return msgs, nil
}
