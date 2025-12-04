package model

type Otp struct {
	Id         uint64 `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	UserId     uint64 `json:"user_id" gorm:"foreignKey:user_id;constraint:OnDelete:CASCADE;not null"`
	EmailOtp   int32  `json:"email_otp"`
	WhatappOtp int32  `json:"whatapp_otp"`
	Status     bool   `json:"status" gorm:"default:0"`
	TimeStamp
}
