package xdb

import "gorm.io/gorm"

type Mail struct {
	gorm.Model

	UserID            int64          `gorm:"not null"`
	SourceType        string         `gorm:"type:varchar(128);not null"`
	SourceUuid        string         `gorm:"type:varchar(255);not null"`
	VisibleType       string         `gorm:"type:varchar(128);not null"`
	ReceiveType       string         `gorm:"type:varchar(128);index:idx_receive_type_uuid;not null"`
	ReceiveUuid       string         `gorm:"type:varchar(255);index:idx_receive_type_uuid"`
	Subject           string         `gorm:"type:text;not null"`
	Body              string         `gorm:"type:text;"`
	Attachments       map[string]any `gorm:"type:json"`
	AttachmentsPicked bool           `gorm:""`
}
