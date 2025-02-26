package xdb

import (
	"encoding/json"
	"errors"

	"github.com/aronfan/plat.mini/xlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	MailVisibleNormal = "normal"
	MailVisibleHidden = "hidden"
)

const (
	MailSourceTypeSystem = "system"
	MailSourceTypeUser   = "user"
)

const (
	MailReceiveTypeOnce    = "once"
	MailReceiveTypeUnlimit = "unlimit"
)

type Mail struct {
	gorm.Model `json:"-"`

	UserID            int64          `json:"user_id" gorm:"not null"`
	VisibleType       string         `json:"visible_type" gorm:"column:visible_type;type:varchar(128);not null"`
	SourceType        string         `json:"source_type" gorm:"type:varchar(128);not null"`
	SourceUuid        string         `json:"source_uuid" gorm:"type:varchar(255);not null"`
	ReceiveType       string         `json:"receive_type" gorm:"column:receive_type;type:varchar(128);index:idx_receive_type_uuid;not null"`
	ReceiveUuid       string         `json:"receive_uuid" gorm:"column:receive_uuid;type:varchar(255);index:idx_receive_type_uuid"`
	Subject           string         `json:"subject" gorm:"type:text;not null"`
	Body              string         `json:"body" gorm:"type:text;"`
	Attachments       map[string]any `json:"attachments" gorm:"-"`
	AttachmentsJSON   string         `json:"-" gorm:"type:json"`
	AttachmentsPicked bool           `json:"-" gorm:""`
}

func InsertMail(ydb *gorm.DB, mail *Mail) (int64, error) {
	if mail.ReceiveType == MailReceiveTypeOnce {
		chk := &Mail{
			ReceiveType: mail.ReceiveType,
			ReceiveUuid: mail.ReceiveUuid,
		}
		query := ydb.Unscoped().Where(chk)
		err := query.First(chk).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, err
			}
		} else {
			xlog.Debug("InsertMail mail already received", zap.String("type", mail.ReceiveType),
				zap.String("uuid", mail.ReceiveUuid))
			return int64(chk.ID), nil
		}
	}

	jsonData, err := json.Marshal(mail.Attachments)
	if err != nil {
		return 0, err
	}
	mail.AttachmentsJSON = string(jsonData)

	if err := ydb.Create(mail).Error; err != nil {
		return 0, err
	}
	return int64(mail.ID), nil
}

func DeleteMail(ydb *gorm.DB, userID, mailID int64) (int64, error) {
	mail := &Mail{
		UserID: userID,
		Model:  gorm.Model{ID: uint(mailID)},
	}
	err := ydb.Delete(mail).Error
	if err != nil {
		return 0, err
	}
	return ydb.RowsAffected, nil
}

func GetMailList(ydb *gorm.DB, userID int64, offset, limit int, seeall bool) ([]*Mail, error) {
	var mails []*Mail
	query := ydb.Where("user_id = ?", userID)
	if !seeall {
		query = query.Where("visible_type = ?", MailVisibleNormal)
	}
	err := query.Offset(offset).Limit(limit).Find(&mails).Error
	if err != nil {
		return nil, err
	}
	return mails, nil
}

func PickMailAttachments(ydb *gorm.DB, userID, mailID int64) error {
	// TODO
	return nil
}
