package output

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/rogafe/go-email/internal/structs"
	"github.com/rogafe/go-email/internal/utils"

	"github.com/emersion/go-message/mail"
)

func WriteImage(eml string, account structs.Account) {
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

		if strings.Contains(p.Header.Get("Content-Type"), "image") {

			// log.Printf("Got attachment==========")

			var folderName string
			if MessageId, err := header.AddressList("Message-Id"); err == nil {
				if len(MessageId) != 0 {
					a := strings.ReplaceAll(MessageId[0].String(), "<", "")
					folderName = strings.ReplaceAll(a, ">", "")
				}
			}
			// imageName := p.Header.Get("name")
			ContentType := p.Header.Get("Content-Type")

			pattern := `name="(.+?)"`

			re := regexp.MustCompile(pattern)
			match := re.FindStringSubmatch(ContentType)

			if len(match) > 1 {
				fmt.Println(match[1])
			} else {
				fmt.Println("No match found")
			}

			imageName := match[1]

			folder := fmt.Sprintf("%s/%s/%s/%s/attachments", account.LocalFolder, account.User, account.RemoteFolder, folderName)
			log.Println(folder)
			utils.CreateFolder(folder)

			log.Printf("Got attachment: %v", imageName)
			b, errp := io.ReadAll(p.Body)
			fmt.Println("errp ===== :", errp)
			err := os.WriteFile(fmt.Sprintf("%s/%s", folder, imageName), b, 0777)
			if err != nil {
				log.Println("attachment err: ", err)
			}
		}

	}

}
