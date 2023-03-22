package output

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
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
	// if MessageId, err := header.AddressList("Message-Id"); err == nil {
	// 	if len(MessageId) != 0 {
	// 		a := strings.ReplaceAll(MessageId[0].String(), "<", "")
	// 		filename = strings.ReplaceAll(a, ">", "")
	// 	}
	// }
	var SenderString, CleanedEmail string
	if Sender, err := header.AddressList("From"); err == nil {
		if len(Sender) != 0 {
			CleanedEmail = strings.ReplaceAll(Sender[0].String(), "[<", "")
			CleanedEmail = strings.ReplaceAll(CleanedEmail, ">]", "")
		}
	}
	//

	re := regexp.MustCompile(`<(.+)>`) // match "<", followed by one or more characters, followed by ">"
	match := re.FindStringSubmatch(CleanedEmail)
	if len(match) > 1 {
		SenderString = strings.Trim(match[1], "<>")
	}

	SubjectString, err := header.Subject()
	if err != nil {
		log.Println(err)
	}
	filename = fmt.Sprintf("%s-%s", SenderString, SubjectString)

	folder := fmt.Sprintf("%s/%s/%s/%s", account.LocalFolder, account.User, account.RemoteFolder, filename)

	utils.CreateFolder(folder)
	err = ioutil.WriteFile(fmt.Sprintf("%s/message.eml", folder), []byte(eml), 0644)
	if err != nil {
		log.Println(err)
	}
}
