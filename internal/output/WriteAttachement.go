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

		switch h := p.Header.(type) {
		case *mail.AttachmentHeader:

			attachmentName, _ := h.Filename()
			var folderName string
			if MessageId, err := header.AddressList("Message-Id"); err == nil {
				if len(MessageId) != 0 {
					a := strings.ReplaceAll(MessageId[0].String(), "<", "")
					folderName = strings.ReplaceAll(a, ">", "")
				}
			}

			folder := fmt.Sprintf("%s/%s/%s/%s", account.LocalFolder, account.User, account.RemoteFolder, folderName)

			utils.CreateFolder(folder)

			b, errp := ioutil.ReadAll(p.Body)
			if errp != nil {
				log.Println(errp)
			}
			log.Println("errp ===== :", errp)
			err := ioutil.WriteFile(fmt.Sprintf("%s/%s", folder, attachmentName), b, 0777)

			if err != nil {
				log.Println("attachment err: ", err)
			}
		}
	}
}
