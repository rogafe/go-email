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

func GetEmails(config structs.Config) {
	var c *client.Client
	var err error

	switch {
	case config.TLS && config.InsecureSkipVerify:
		c, err = client.DialTLS(config.Uri, &tls.Config{InsecureSkipVerify: true})
	case config.TLS && !config.InsecureSkipVerify:
		c, err = client.DialTLS(config.Uri, nil)
	case !config.TLS:
		c, err = client.Dial(config.Uri)
	}

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	//LOGIN
	switch config.Oauth2 {
	case "gmail":
		token := auth.GoogleOauth()

		err = c.Authenticate(sasl.NewOAuthBearerClient(&sasl.OAuthBearerOptions{
			Username: config.User,
			Token:    token.AccessToken,
		}))
		if err != nil {
			log.Println(err)
		}
		log.Println("Logged in")
	default:
		if err := c.Login(config.User, config.Password); err != nil {
			log.Fatal(err)
		}
		log.Println("Logged in")
	}

	// Select INBOX
	mbox, err := c.Select(config.RemoteFolder, false)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Downloding %d emails\n", mbox.Messages)
	seqset := new(imap.SeqSet)
	seqset.AddRange(1, mbox.Messages)
	// seqset.AddRange(mbox.Messages, mbox.Messages)
	// seqset.AddRange(mbox.Messages-5, mbox.Messages)
	// var Messages []*imap.Message
	messages := make(chan *imap.Message, mbox.Messages)
	done := make(chan error, 1)
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	go func() {
		done <- c.Fetch(seqset, items, messages)
	}()

	log.Println("All the e-mail have been downloaded, converting to EML")

	sl := utils.ChanToSlice(messages).([]*imap.Message)

	for i, msg := range sl {

		log.Printf("Email %d out of %d", i, mbox.Messages)
		if msg == nil {
			log.Fatal("Server didn't returned message")
		}
		r := msg.GetBody(&section)
		if r == nil {
			log.Fatal("Server didn't returned message body")
		}
		eml := utils.StreamToString(r)

		for _, out := range config.OutputTypes {
			switch out {
			case "eml":
				go output.WriteEML(eml, config)
			case "html":
				go output.WriteHTML(eml, config)
			case "json":
				go output.WriteJSON(eml, config)
			case "attachement":
				log.Println(out)
				go output.WriteAttachement(eml, config)

			}
		}
	}

}
