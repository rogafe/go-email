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

func WriteHTML(eml string, account structs.Account) {
	mr, err := mail.CreateReader(strings.NewReader(eml))
	if err != nil {
		log.Println(err)
	}

	// Print some info about the message
	header := mr.Header
	var Body []byte
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		log.Println(p.Header.Get("Content-Type"))

		if strings.Contains(p.Header.Get("Content-Type"), "text/html") {
			b, err := ioutil.ReadAll(p.Body)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(len(string(b)))
			Body = b

		} else {
			log.Println(p.Header.Get("Content-Type"))
		}
	}
	// log.Println(len(string(Body)))

	var filename string
	if MessageId, err := header.AddressList("Message-Id"); err == nil {
		if len(MessageId) != 0 {
			a := strings.ReplaceAll(MessageId[0].String(), "<", "")
			filename = strings.ReplaceAll(a, ">", "")
		}
	}

	folder := fmt.Sprintf("%s/%s/%s/%s", account.LocalFolder, account.User, account.RemoteFolder, filename)

	utils.CreateFolder(folder)
	err = ioutil.WriteFile(fmt.Sprintf("%s/message.html", folder), []byte(Body), 0644)
	if err != nil {
		log.Println(err)
	}

}
