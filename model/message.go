package model

import "gorm.io/gorm"

type GroupMessage struct {
	gorm.Model

	SenderID uint `json:"uid" gorm:"column:uid"`
	Sender   User `json:"sender" gorm:"foreignKey:SenderID"`

	Content string `json:"content" gorm:"type:text"`
}

type Message struct {
	gorm.Model

	SID    uint `json:"sid" gorm:"column:sid"`
	Sender User `json:"sender" gorm:"foreignKey:SID"`

	RID      uint `json:"rid" gorm:"column:rid"`
	Receiver User `json:"receiver" gorm:"foreignKey:RID"`

	Type    int    `json:"type"`
	Content []byte `json:"content" gorm:""`
}
