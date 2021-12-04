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

func WriteHTML(eml string, account structs.Account, outputType string) (HtmlString string) {
	mr, err := mail.CreateReader(strings.NewReader(eml))
	if err != nil {
		log.Println(err)
	}

	header := mr.Header
	var Body []byte
	var ImageStruct []structs.Image
	for {
		var tmpImage structs.Image
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if strings.Contains(p.Header.Get("Content-Type"), "text/html") {
			b, err := ioutil.ReadAll(p.Body)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(len(string(b)))
			Body = b

		} else if strings.Contains(p.Header.Get("Content-Type"), "image") {

			tmpImage.ImageType = p.Header.Get("Content-Type")
			tmpImage.ImageContentID = p.Header.Get("Content-ID")
			tmpImage.ImageName = p.Header.Get("Content-Description")

			log.Println(p.Header.Get("filename"))

			ImageStruct = append(ImageStruct, tmpImage)
		} else {
			log.Println(p.Header.Get("Content-Type"))
		}
	}

	HtmlString = string(Body)

	// clean charset iso-8859-1

	HtmlString = strings.Replace(HtmlString, `<meta http-equiv="Content-Type" content="text/html; charset=iso-8859-1">`, `<meta http-equiv="Content-Type" content="text/html; charset=utf-8">`, -1)

	var filename string
	if MessageId, err := header.AddressList("Message-Id"); err == nil {
		if len(MessageId) != 0 {
			a := strings.ReplaceAll(MessageId[0].String(), "<", "")
			filename = strings.ReplaceAll(a, ">", "")
		}
	}
	folder := fmt.Sprintf("%s/%s/%s/%s", account.LocalFolder, account.User, account.RemoteFolder, filename)

	for _, IMG := range ImageStruct {
		log.Println(IMG.ImageContentID)
		cid := fmt.Sprintf("cid:%s", IMG.ImageContentID)
		a := strings.ReplaceAll(cid, "<", "")
		cid = strings.ReplaceAll(a, ">", "")

		if strings.Contains(HtmlString, cid) {
			// 	log.Panic("tolo")
			log.Println("yolo")
			HtmlString = strings.Replace(HtmlString, cid, IMG.ImageName, -1)
		}
	}

	switch outputType {
	case "file":
		utils.CreateFolder(folder)
		err = ioutil.WriteFile(fmt.Sprintf("%s/message.html", folder), []byte(HtmlString), 0644)
		if err != nil {
			log.Println(err)
		}
	case "string":
		return HtmlString
	}

	return ""

}
