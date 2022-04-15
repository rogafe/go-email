package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/emersion/go-message/charset"
	"github.com/rogafe/go-email/internal/email"
	"github.com/rogafe/go-email/internal/structs"
	"gopkg.in/ini.v1"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("go-email: ")
	log.SetOutput(os.Stderr)
}

func main() {
	cfg, err := ini.Load("config.cfg")
	if err != nil {
		log.Panicf("Fail to read file: %v", err)

	}

	config := structs.Config{
		LocalFolder: cfg.Section("go-email").Key("local_folder").String(),
		Wkhtmltopdf: cfg.Section("go-email").Key("wkhtmltopdf").String(),
	}

	for _, section := range cfg.SectionStrings() {
		log.Println(section)
		if section != ini.DefaultSection && section != "go-email" {

			Name := section
			InsecureSkipVerify, _ := cfg.Section(Name).Key("insecureskipverify").Bool()
			TLS, _ := cfg.Section(Name).Key("ssl").Bool()

			account := structs.Account{
				Name:               Name,
				LocalFolder:        config.LocalFolder,
				Uri:                fmt.Sprintf("%s:%s", cfg.Section(Name).Key("host").String(), cfg.Section(Name).Key("port").String()),
				User:               cfg.Section(Name).Key("username").String(),
				Password:           cfg.Section(Name).Key("password").String(),
				Oauth2:             cfg.Section(Name).Key("oauth2").String(),
				RemoteFolder:       cfg.Section(Name).Key("remote_folder").String(),
				TLS:                TLS,
				InsecureSkipVerify: InsecureSkipVerify,
				OutputTypes:        strings.Split(cfg.Section("go-email").Key("output_types").String(), ","),
			}

			config.Accounts = append(config.Accounts, account)

		}

	}

	email.Get(config)

}
