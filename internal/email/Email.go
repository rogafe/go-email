package email

import (
	"log"

	"github.com/rogafe/go-email/internal/structs"
)

func Get(config structs.Config) {
	for _, account := range config.Accounts {
		log.Println(account.Name)
		if account.RemoteFolder == "all" {
			GetAllEmails(account)
		} else {
			GetEmails(account)
		}
	}
}
