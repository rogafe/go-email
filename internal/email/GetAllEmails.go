package email

import (
	"crypto/tls"
	"log"

	"github.com/rogafe/go-email/internal/auth"
	"github.com/rogafe/go-email/internal/output"
	"github.com/rogafe/go-email/internal/structs"
	"github.com/rogafe/go-email/internal/utils"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
)

func GetAllEmails(account structs.Account) {
	var c *client.Client
	var err error

	switch {
	case account.TLS && account.InsecureSkipVerify:
		c, err = client.DialTLS(account.Uri, &tls.Config{InsecureSkipVerify: true})
	case account.TLS && !account.InsecureSkipVerify:
		c, err = client.DialTLS(account.Uri, nil)
	case !account.TLS:
		c, err = client.Dial(account.Uri)
	}

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	switch account.Oauth2 {
	case "gmail":
		token := auth.GoogleOauth(account)

		err = c.Authenticate(sasl.NewOAuthBearerClient(&sasl.OAuthBearerOptions{
			Username: account.User,
			Token:    token.AccessToken,
		}))
		if err != nil {
			log.Println(err)
		}
		log.Println("Logged in")
	default:
		if err := c.Login(account.User, account.Password); err != nil {
			log.Fatal(err)
		}
		log.Println("Logged in")
	}

	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	var folders []string
	for m := range mailboxes {
		folders = append(folders, m.Name)
	}

	for _, folder := range folders {

		// Select INBOX
		mbox, err := c.Select(folder, false)
		if err != nil {
			log.Fatal(err)
		}

		if mbox.Messages != 0 {

			log.Printf("Downloding %d emails\n", mbox.Messages)
			seqset := new(imap.SeqSet)
			seqset.AddRange(1, mbox.Messages)

			messages := make(chan *imap.Message, mbox.Messages)
			done := make(chan error, 1)
			var section imap.BodySectionName
			items := []imap.FetchItem{section.FetchItem()}

			go func() {
				done <- c.Fetch(seqset, items, messages)
			}()

			log.Println("All the e-mail have been downloaded, converting to EML")

			// sl := utils.ChanToSlice(messages).([]*imap.Message)
			var i int
			for msg := range messages {

				log.Printf("Email %d out of %d", i, mbox.Messages)
				if msg == nil {
					log.Fatal("Server didn't returned message")
				}
				r := msg.GetBody(&section)
				if r == nil {
					log.Fatal("Server didn't returned message body")
				}
				eml := utils.StreamToString(r)
				account.RemoteFolder = folder

				log.Println(account.OutputTypes)
				for _, out := range account.OutputTypes {
					log.Println(out)
					switch out {
					case "eml":
						log.Println(out)
						go output.WriteEML(eml, account)
					case "html":
						log.Println(out)
						go output.WriteHTML(eml, account)
					case "json":
						log.Println(out)
						go output.WriteJSON(eml, account)
					case "attachement":
						log.Println(out)
						go output.WriteAttachement(eml, account)

					}
				}
				i++

			}
		}
	}

}
