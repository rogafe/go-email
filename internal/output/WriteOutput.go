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
				// 		// go WriteEML(eml, account)
				WriteEML(msg.EML, msg.Account)
			case "html":
				log.Println("html")
				WriteHTML(msg.EML, msg.Account, "file")
			case "json":
				log.Println("json")
				// go WriteJSON(eml, account)
				WriteJSON(msg.EML, msg.Account)
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
	// for _, out := range account.OutputTypes {
	// 	log.Println(out)
	// 	// 	switch out {
	// 	// 	case "eml":
	// 	// 		log.Println("eml")
	// 	// 		// go WriteEML(eml, account)
	// 	// 		WriteEML(eml, account)
	// 	// 	case "html":
	// 	// 		log.Println("html")
	// 	// 		// go WriteHTML(eml, account, "file")
	// 	// 		WriteHTML(eml, account, "file")
	// 	// 	case "json":
	// 	// 		log.Println("json")
	// 	// 		// go WriteJSON(eml, account)
	// 	// 		WriteJSON(eml, account)
	// 	// 	case "attachement":
	// 	// 		log.Println("attachement")
	// 	// 		// go WriteAttachement(eml, account)
	// 	// 		WriteAttachement(eml, account)
	// 	// 	case "image":
	// 	// 		log.Println("image")
	// 	// 		// go WriteImage(eml, account)
	// 	// 		WriteImage(eml, account)
	// 	// 	case "pdf":
	// 	// 		log.Println("pdf")
	// 	// 		// go WritePDF(eml, account)
	// 	// 		WritePDF(eml, account)
	// 	// 	}
	// }
}
