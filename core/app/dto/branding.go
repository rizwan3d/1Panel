package dto

// LoginBranding contains the public visual branding shown before authentication.
type LoginBranding struct {
	WelcomeMessage string `json:"welcomeMessage"`
	Logo           string `json:"logo"`
	WebsiteIcon    string `json:"websiteIcon"`
}

// LoginBrandingUpdate replaces the editable login branding configuration.
type LoginBrandingUpdate struct {
	WelcomeMessage string `json:"welcomeMessage"`
	Logo           string `json:"logo"`
	WebsiteIcon    string `json:"websiteIcon"`
}
