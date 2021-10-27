package output

import (
	"fmt"
	"go-email/internal/structs"
	"go-email/internal/utils"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
)

func WriteEML(eml string, config structs.Config) {
	mr, err := mail.CreateReader(strings.NewReader(eml))
	if err != nil {
		log.Println(err)
	}

	header := mr.Header
	var filename string
	if MessageId, err := header.AddressList("Message-Id"); err == nil {
		a := strings.ReplaceAll(MessageId[0].String(), "<", "")
		filename = strings.ReplaceAll(a, ">", "")
	}

	folder := fmt.Sprintf("%s/%s", config.LocalFolder, filename)

	utils.CreateFolder(folder)
	err = ioutil.WriteFile(fmt.Sprintf("%s/message.eml", folder), []byte(eml), 0644)
	if err != nil {
		log.Println(err)
	}
}
