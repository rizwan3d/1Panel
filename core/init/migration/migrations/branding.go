package migrations

import (
	"errors"

	"github.com/1Panel-dev/1Panel/core/app/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddLoginBrandingSettings = &gormigrate.Migration{
	ID: "20260819-add-login-branding-settings",
	Migrate: func(tx *gorm.DB) error {
		addIfMissing := func(key, value string) error {
			var setting model.Setting
			if err := tx.Where("key = ?", key).First(&setting).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return tx.Create(&model.Setting{Key: key, Value: value}).Error
				}
				return err
			}
			return nil
		}

		for _, item := range []struct {
			key   string
			value string
		}{
			{key: "LoginWelcomeMessage", value: ""},
			{key: "LoginLogo", value: ""},
			{key: "WebsiteIcon", value: ""},
		} {
			if err := addIfMissing(item.key, item.value); err != nil {
				return err
			}
		}
		return nil
	},
}
