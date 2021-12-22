package output

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/rogafe/go-email/internal/structs"
	"github.com/rogafe/go-email/internal/utils"

	"github.com/emersion/go-message/mail"
)

func WriteAttachement(eml string, account structs.Account) {
	mr, err := mail.CreateReader(strings.NewReader(eml))
	if err != nil {
		log.Println(err)
	}

	header := mr.Header
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// if strings.Contains(p.Header.Get("Content-Type"), "image") {

		// 	// log.Printf("Got attachment==========")

		// 	log.Println(p.Header.Get("Content-Description"))
		// }

		switch h := p.Header.(type) {
		case *mail.AttachmentHeader:

			log.Printf("Got attachment==========")
			// This is an attachment
			attachmentName, _ := h.Filename()
			// log.Panic("39: " + attachmentName)
			var folderName string
			if MessageId, err := header.AddressList("Message-Id"); err == nil {
				if len(MessageId) != 0 {
					a := strings.ReplaceAll(MessageId[0].String(), "<", "")
					folderName = strings.ReplaceAll(a, ">", "")
				}
			}

			folder := fmt.Sprintf("%s/%s/%s/%s", account.LocalFolder, account.User, account.RemoteFolder, folderName)

			utils.CreateFolder(folder)

			log.Printf("Got attachment: %v", attachmentName)
			b, errp := ioutil.ReadAll(p.Body)
			fmt.Println("errp ===== :", errp)
			err := ioutil.WriteFile(fmt.Sprintf("%s/%s", folder, attachmentName), b, 0777)

			if err != nil {
				log.Println("attachment err: ", err)
			}
		}
	}
}
