package email

import (
	"crypto/tls"
	"log"

	pb "github.com/cheggaaa/pb/v3"
	"github.com/rogafe/go-email/internal/auth"
	"github.com/rogafe/go-email/internal/output"
	"github.com/rogafe/go-email/internal/structs"
	"github.com/rogafe/go-email/internal/utils"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
)

func GetEmails(account structs.Account) {
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

	//LOGIN
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

	// Select INBOX
	mbox, err := c.Select(account.RemoteFolder, false)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Downloding %d emails\n", mbox.Messages)
	seqset := new(imap.SeqSet)

	//set range to all message in inbox might want to control that in the future
	seqset.AddRange(1, mbox.Messages)

	messages := make(chan *imap.Message, mbox.Messages)
	done := make(chan error, 1)
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	go func() {
		done <- c.Fetch(seqset, items, messages)
	}()

	var i int

	var Messages []structs.Messages
	bar := pb.Full.Start(int(mbox.Messages))
	for msg := range messages {

		bar.Increment()
		if msg == nil {
			log.Fatal("Server didn't returned message")
		}
		r := msg.GetBody(&section)
		if r == nil {
			log.Fatal("Server didn't returned message body")
		}
		eml := utils.StreamToString(r)
		msgs := structs.Messages{
			EML:     eml,
			Account: account,
		}
		Messages = append(Messages, msgs)

		i++
	}
	bar.Finish()

	log.Println("All the e-mail have been downloaded")
	if err := <-done; err != nil {
		log.Fatal(err)
	}

	output.WriteOutput(Messages)

}
