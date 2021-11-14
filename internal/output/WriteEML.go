package output

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/rogafe/go-email/internal/structs"
	"github.com/rogafe/go-email/internal/utils"

	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
)

func WriteEML(eml string, account structs.Account) {
	mr, err := mail.CreateReader(strings.NewReader(eml))
	if err != nil {
		log.Println(err)
	}

	header := mr.Header
	var filename string
	if MessageId, err := header.AddressList("Message-Id"); err == nil {
		if len(MessageId) != 0 {
			a := strings.ReplaceAll(MessageId[0].String(), "<", "")
			filename = strings.ReplaceAll(a, ">", "")
		}
	}

	folder := fmt.Sprintf("%s/%s/%s/%s", account.LocalFolder, account.User, account.RemoteFolder, filename)

	utils.CreateFolder(folder)
	err = ioutil.WriteFile(fmt.Sprintf("%s/message.eml", folder), []byte(eml), 0644)
	if err != nil {
		log.Println(err)
	}
}
