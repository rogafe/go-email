package output

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/rogafe/go-email/internal/structs"
	"github.com/rogafe/go-email/internal/utils"

	"github.com/emersion/go-message/mail"
)

func WriteJSON(eml string, config structs.Config) {
	log.Println("json")

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

	// var Test []string
	// // Process each message's part
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		// switch type := p.Header.Get("Content-Type") {
		// case strings.Contains(type, "text/html"):
		// 	b, _ := ioutil.ReadAll(p.Body)
		// 	Email.Body = string(b)
		// 	log.Println(len(Email.Body))
		// }
		if strings.Contains(p.Header.Get("Content-Type"), "text/html") {
			b, _ := ioutil.ReadAll(p.Body)
			Email.Body = string(b)
		} else {
			log.Println(p.Header.Get("Content-Type"))
		}
	}

	// log.Println(len(Email.Body))

	json, err := json.MarshalIndent(Email, "", " ")
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(string(json))

	var filename string
	if MessageId, err := header.AddressList("Message-Id"); err == nil {
		if len(MessageId) != 0 {
			Email.ID = MessageId[0].String()
			a := strings.ReplaceAll(MessageId[0].String(), "<", "")
			filename = strings.ReplaceAll(a, ">", "")
		}
	}

	folder := fmt.Sprintf("%s/%s/%s", config.LocalFolder, config.RemoteFolder, filename)

	utils.CreateFolder(folder)
	err = ioutil.WriteFile(fmt.Sprintf("%s/message.json", folder), json, 0644)
	if err != nil {
		log.Println(err)
	}

}
