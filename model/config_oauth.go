package model

type OAuthProviderSettings struct {
	Enable          *bool
	Secret          *string
	TrustedDomain   *string
	EmailField      *string
	NameField       *string
	Id              *string
	Scope           *string
	AuthEndpoint    *string
	TokenEndpoint   *string
	UserApiEndpoint *string
}

func (o *Config) GetOAuthProvider(service string) *OAuthProviderSettings {
	a := o.OAuthSettings[service]
	return &a
}
