package structs

type Config struct {
	Uri                string
	User               string
	Password           string
	RemoteFolder       string
	TLS                bool
	InsecureSkipVerify bool
	LocalFolder        string
	OutputTypes        []string
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
