package output

import (
	"log"

	"github.com/rogafe/go-email/internal/structs"
)

func WriteOutput(eml string, account structs.Account) {
	for _, out := range account.OutputTypes {
		log.Println(out)
		switch out {
		case "eml":
			log.Println(out)
			go WriteEML(eml, account)
		case "html":
			log.Println(out)
			go WriteHTML(eml, account, "file")
		case "json":
			log.Println(out)
			go WriteJSON(eml, account)
		case "attachement":
			log.Println(out)
			go WriteAttachement(eml, account)

		case "image":
			log.Println(out)
			go WriteImage(eml, account)
		case "pdf":
			log.Println(out)
			go WritePDF(eml, account)
		}
	}
}
