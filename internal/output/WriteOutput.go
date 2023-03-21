package output

import (
	"log"

	"github.com/rogafe/go-email/internal/structs"
)

func WriteOutput(Messages []structs.Messages) {
	// log.Println(account)

	for _, msg := range Messages {
		for _, outType := range msg.Account.OutputTypes {
			switch outType {
			case "eml":
				log.Println("eml")
				WriteEML(msg.EML, msg.Account)
			case "html":
				log.Println("html")
				WriteHTML(msg.EML, msg.Account, "file")
			case "json":
				log.Println("json")
				// go WriteJSON(eml, account)
				// WriteJSON(msg.EML, msg.Account)
			case "attachement":
				log.Println("attachement")
				// go WriteAttachement(eml, account)
				WriteAttachement(msg.EML, msg.Account)
			case "image":
				log.Println("image")
				// go WriteImage(eml, account)
				WriteImage(msg.EML, msg.Account)
			case "pdf":
				log.Println("pdf")
				// go WritePDF(eml, account)
				WritePDF(msg.EML, msg.Account)

			}
		}
	}
}
