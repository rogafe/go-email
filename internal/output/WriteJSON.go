package output

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/rogafe/go-email/internal/structs"
	"github.com/rogafe/go-email/internal/utils"

	"github.com/emersion/go-message/mail"
)

func WriteJSON(eml string, account structs.Account) {
	mr, err := mail.CreateReader(strings.NewReader(eml))
	if err != nil {
		log.Println(err)
	}

	// Print some info about the message
	header := mr.Header
	var Email structs.Email

	if date, err := header.Date(); err == nil {
		Email.Date = date.String()
	}
	if from, err := header.AddressList("From"); err == nil {
		for _, f := range from {
			Email.From = append(Email.From, f.String())
		}
	}
	if to, err := header.AddressList("To"); err == nil {
		for _, t := range to {
			Email.To = append(Email.To, t.String())

		}
	}
	if subject, err := header.Subject(); err == nil {
		Email.Subject = subject
	}

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if strings.Contains(p.Header.Get("Content-Type"), "text/html") {
			b, _ := ioutil.ReadAll(p.Body)
			Email.Body = string(b)
		} else {
			log.Println(p.Header.Get("Content-Type"))
		}
	}

	json, err := json.MarshalIndent(Email, "", " ")
	if err != nil {
		log.Println(err)
	}

	var filename string
	// if MessageId, err := header.AddressList("Message-Id"); err == nil {
	// 	if len(MessageId) != 0 {
	// 		Email.ID = MessageId[0].String()
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
	err = ioutil.WriteFile(fmt.Sprintf("%s/message.json", folder), json, 0644)
	if err != nil {
		log.Println(err)
	}

}
