package model

import "gorm.io/gorm"

type GroupMessage struct {
	gorm.Model

	SID     string `json:"sid" gorm:"size:255"`
	RID     string `json:"rid" gorm:"size:255"`
	Type    int    `json:"type"`
	Content string `json:"content" gorm:"type:text"`
}

type Message struct {
	gorm.Model

	SID    uint `json:"sid" gorm:"column:sid"`
	Sender User `json:"sender" gorm:"foreignKey:SID"`

	RID      uint `json:"rid" gorm:"column:rid"`
	Receiver User `json:"receiver" gorm:"foreignKey:RID"`

	Type    int    `json:"type"`
	Content string `json:"content" gorm:""`
}
