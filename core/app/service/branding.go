package service

import (
	"encoding/base64"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/1Panel-dev/1Panel/core/app/dto"
	"github.com/1Panel-dev/1Panel/core/app/repo"
)

const (
	maxWelcomeMessageLength = 200
	maxLogoBytes            = 2 * 1024 * 1024
	maxWebsiteIconBytes     = 1 * 1024 * 1024
)

var logoMimeTypes = map[string]struct{}{
	"image/png":  {},
	"image/jpeg": {},
	"image/webp": {},
}

var websiteIconMimeTypes = map[string]struct{}{
	"image/png":                {},
	"image/webp":               {},
	"image/x-icon":             {},
	"image/vnd.microsoft.icon": {},
}

func (u *SettingService) GetLoginBranding() (*dto.LoginBranding, error) {
	welcome, err := settingRepo.Get(repo.WithByKey("LoginWelcomeMessage"))
	if err != nil {
		return nil, err
	}
	logo, err := settingRepo.Get(repo.WithByKey("LoginLogo"))
	if err != nil {
		return nil, err
	}
	websiteIcon, err := settingRepo.Get(repo.WithByKey("WebsiteIcon"))
	if err != nil {
		return nil, err
	}
	return &dto.LoginBranding{
		WelcomeMessage: welcome.Value,
		Logo:           logo.Value,
		WebsiteIcon:    websiteIcon.Value,
	}, nil
}

func (u *SettingService) UpdateLoginBranding(req dto.LoginBrandingUpdate) error {
	if utf8.RuneCountInString(req.WelcomeMessage) > maxWelcomeMessageLength {
		return errors.New("welcome message cannot exceed 200 characters")
	}
	if !validBrandingImage(req.Logo, logoMimeTypes, maxLogoBytes) {
		return errors.New("invalid logo image")
	}
	if !validBrandingImage(req.WebsiteIcon, websiteIconMimeTypes, maxWebsiteIconBytes) {
		return errors.New("invalid website icon")
	}
	if err := settingRepo.Update("LoginWelcomeMessage", strings.TrimSpace(req.WelcomeMessage)); err != nil {
		return err
	}
	if err := settingRepo.Update("LoginLogo", req.Logo); err != nil {
		return err
	}
	return settingRepo.Update("WebsiteIcon", req.WebsiteIcon)
}

func validBrandingImage(value string, allowedMimeTypes map[string]struct{}, maxBytes int) bool {
	if value == "" {
		return true
	}
	parts := strings.SplitN(value, ",", 2)
	if len(parts) != 2 || !strings.HasPrefix(parts[0], "data:") || !strings.HasSuffix(parts[0], ";base64") {
		return false
	}
	mimeType := strings.TrimSuffix(strings.TrimPrefix(parts[0], "data:"), ";base64")
	if _, ok := allowedMimeTypes[mimeType]; !ok {
		return false
	}
	if len(parts[1]) > ((maxBytes+2)/3)*4+4 {
		return false
	}
	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	return err == nil && len(decoded) <= maxBytes
}
