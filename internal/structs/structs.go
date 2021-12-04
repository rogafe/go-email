package structs

type Config struct {
	Days        int64     `json:"days"`
	LocalFolder string    `json:"local_folder"`
	Wkhtmltopdf string    `json:"wkhtmltopdf"`
	Accounts    []Account `json:"accounts"`
	Verbose     bool      `json:"verbose"`
}

type Account struct {
	Name               string   `json:"name"`
	LocalFolder        string   `json:"local_folder"`
	RemoteFolder       string   `json:"remote_folder"`
	Port               string   `json:"port"`
	TLS                bool     `json:"TLS"`
	InsecureSkipVerify bool     `json:"insecure_skip_verify"`
	Uri                string   `json:"uri"`
	User               string   `json:"username"`
	Password           string   `json:"password"`
	Oauth2             string   `json:"oauth2"`
	OutputTypes        []string `json:"output_type"`
}

// Generated by https://quicktype.io

type Email struct {
	ID          string        `json:"Id"`
	Subject     string        `json:"Subject"`
	From        []string      `json:"From"`
	To          []string      `json:"To"`
	Cc          []interface{} `json:"Cc"`
	Date        string        `json:"Date"`
	UTC         string        `json:"Utc"`
	Attachments []interface{} `json:"Attachments"`
	WithHTML    bool          `json:"WithHtml"`
	WithText    bool          `json:"WithText"`
	Body        string        `json:"Body"`
}

// Generated by https://quicktype.io

type GoogleToken struct {
	Installed Installed `json:"installed"`
}

type Installed struct {
	ClientID                string   `json:"client_id"`
	ProjectID               string   `json:"project_id"`
	AuthURI                 string   `json:"auth_uri"`
	TokenURI                string   `json:"token_uri"`
	AuthProviderX509CERTURL string   `json:"auth_provider_x509_cert_url"`
	ClientSecret            string   `json:"client_secret"`
	RedirectUris            []string `json:"redirect_uris"`
}

type Image struct {
	ImageType      string
	ImageContentID string
	ImageName      string
}
