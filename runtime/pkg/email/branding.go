package email

import "fmt"

// Branding holds configurable brand values for email templates.
// The zero value is not usable; call DefaultBranding for Rill defaults.
type Branding struct {
	LogoURL         string // header logo image
	CompanyName     string // copyright line: "Rill Data, Inc."
	ProductName     string // product references: "Rill"
	PrimaryColor    string // brand color for buttons/links: "#4736F5"
	SupportEmail    string // support contact: "support@rilldata.com"
	Address         string // legal footer address
	ContactURL      string // "Contact us" link
	CommunityURL    string // "Community" link
	PrivacyURL      string // "Privacy Policy" link
	DocsURL         string // documentation root
	ReleaseNotesURL string // release notes page
	ChatURL         string // in-app chat / support chat link
	TeamSignature   string // sign-off: "The Rill Team"
}

// DefaultBranding returns Rill's default branding values.
func DefaultBranding() *Branding {
	return &Branding{
		LogoURL:         "https://cdn.rilldata.com/email-transactional/rill-logo-purple-square.png",
		CompanyName:     "Rill Data, Inc.",
		ProductName:     "Rill",
		PrimaryColor:    "#4736F5",
		SupportEmail:    "support@rilldata.com",
		Address:         "18 Bartol St. • San Francisco • CA",
		ContactURL:      "https://www.rilldata.com/contact",
		CommunityURL:    "https://bit.ly/3unvA05",
		PrivacyURL:      "https://www.rilldata.com/legal/privacy",
		DocsURL:         "https://docs.rilldata.com",
		ReleaseNotesURL: "https://docs.rilldata.com/notes",
		ChatURL:         "https://docs.rilldata.com/contact#in-app-chat",
		TeamSignature:   "The Rill Team",
	}
}

// Option configures the email Client.
type Option func(*Client)

// WithBranding sets custom branding for email templates.
func WithBranding(b *Branding) Option {
	return func(c *Client) {
		if b != nil {
			c.branding = b
		}
	}
}

// BrandingFromOrg returns a *Branding that overlays an organization's display
// name and logo URL on top of the Rill defaults.  Empty/zero fields keep the
// default.  logoURL should already be a resolved URL (not an asset ID).
func BrandingFromOrg(displayName, logoURL string) *Branding {
	b := DefaultBranding()
	if displayName != "" {
		b.CompanyName = displayName
		b.ProductName = displayName
		b.TeamSignature = fmt.Sprintf("The %s Team", displayName)
	}
	if logoURL != "" {
		b.LogoURL = logoURL
	}
	return b
}
