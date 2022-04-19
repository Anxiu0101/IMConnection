package model

import "gorm.io/gorm"

type GroupMessage struct {
	gorm.Model

	SenderID uint `json:"uid" gorm:"column:uid"`
	Sender   User `json:"sender" gorm:"foreignKey:SenderID"`

	Content string `json:"content" gorm:"type:text"`
}

type Message struct {
}
