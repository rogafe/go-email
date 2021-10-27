package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"go-email/internal/email"
	"go-email/internal/structs"

	_ "github.com/emersion/go-message/charset"
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

	InsecureSkipVerify, _ := cfg.Section("email").Key("insecureskipverify").Bool()
	TLS, _ := cfg.Section("email").Key("ssl").Bool()

	config := structs.Config{
		Uri:                fmt.Sprintf("%s:%s", cfg.Section("email").Key("host").String(), cfg.Section("email").Key("port").String()),
		User:               cfg.Section("email").Key("username").String(),
		Password:           cfg.Section("email").Key("password").String(),
		RemoteFolder:       cfg.Section("email").Key("remote_folder").String(),
		TLS:                TLS,
		InsecureSkipVerify: InsecureSkipVerify,
		LocalFolder:        cfg.Section("go-email").Key("local_folder").String(),
		OutputTypes:        strings.Split(cfg.Section("go-email").Key("output_types").String(), ","),
	}
	if config.RemoteFolder == "all" {
		email.GetAllEmails(config)
	} else {
		email.GetEmails(config)
	}
}
