package output

import (
	"fmt"
	"log"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/emersion/go-message/mail"
	"github.com/rogafe/go-email/internal/structs"
	"github.com/rogafe/go-email/internal/utils"
)

func WritePDF(eml string, account structs.Account) {

	mr, err := mail.CreateReader(strings.NewReader(eml))
	if err != nil {
		log.Panicln(err)
	}

	header := mr.Header

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Panic(err)
	}
	pdfg.Dpi.Set(600)
	pdfg.NoCollate.Set(true)
	pdfg.MarginBottom.Set(40)
	pdfg.Cover.EnableLocalFileAccess.Set(true)
	pdfg.Grayscale.Set(true)

	file := "message.pdf"
	var filename string
	if MessageId, err := header.AddressList("Message-Id"); err == nil {
		if len(MessageId) != 0 {
			a := strings.ReplaceAll(MessageId[0].String(), "<", "")
			filename = strings.ReplaceAll(a, ">", "")
		}
	}
	folder := fmt.Sprintf("%s/%s/%s/%s", account.LocalFolder, account.User, account.RemoteFolder, filename)

	out := WriteHTML(eml, account, "string")

	page := wkhtmltopdf.NewPageReader(strings.NewReader(out))
	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Panicln(err)
	}

	utils.CreateFolder(folder)

	err = pdfg.WriteFile(fmt.Sprintf("%s/%s", folder, file))
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println("Done")

}
