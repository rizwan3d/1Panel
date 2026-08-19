package v2

import (
	"github.com/1Panel-dev/1Panel/core/app/api/v2/helper"
	"github.com/1Panel-dev/1Panel/core/app/dto"
	"github.com/1Panel-dev/1Panel/core/app/service"
	"github.com/gin-gonic/gin"
)

var loginBrandingService = &service.SettingService{}

// GetLoginBranding returns the public branding used by the login page.
func (b *BaseApi) GetLoginBranding(c *gin.Context) {
	branding, err := loginBrandingService.GetLoginBranding()
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, branding)
}

// UpdateLoginBranding updates the login welcome message, logo, and website icon.
func (b *BaseApi) UpdateLoginBranding(c *gin.Context) {
	var req dto.LoginBrandingUpdate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := loginBrandingService.UpdateLoginBranding(req); err != nil {
		helper.BadRequest(c, err)
		return
	}
	helper.Success(c)
}
