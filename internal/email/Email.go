package email

import (
	"github.com/rogafe/go-email/internal/structs"
)

func Get(config structs.Config) {
	for _, account := range config.Accounts {
		if account.RemoteFolder == "all" {
			GetAllEmails(account)
		} else {
			GetEmails(account)
		}
	}
}
