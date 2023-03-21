package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/rogafe/go-email/internal/server"
	"github.com/rogafe/go-email/internal/structs"
	"github.com/rogafe/go-email/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func LoadConfiguration(file string) (account *oauth2.Config) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	account, err = google.ConfigFromJSON(b, "https://mail.google.com")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to account: %v", err)
	}

	return account
}

func LoadToken(file string) (token *oauth2.Token) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	err = json.Unmarshal(b, &token)
	if err != nil {
		log.Println(err)
	}
	return token
}

func WriteToken(token *oauth2.Token, account structs.Account) {
	tokenStr := fmt.Sprintf("./%s/%s/token.json", account.LocalFolder, account.User)

	tokenJson, err := json.MarshalIndent(token, "", "\t")
	if err != nil {
		log.Println(err)
	}

	utils.CreateFolder(fmt.Sprintf("./%s/%s", account.LocalFolder, account.User))
	err = os.WriteFile(tokenStr, tokenJson, 0666)
	if err != nil {
		log.Println(err)
	}
}

func GoogleOauth(account structs.Account) (token *oauth2.Token) {
	tokenStr := fmt.Sprintf("./%s/%s/token.json", account.LocalFolder, account.User)

	exist, err := utils.FileExists(tokenStr)
	if err != nil {
		log.Println(err)
	}

	if !exist {
		token = GetToken(account)
	} else {
		token = LoadToken(fmt.Sprintf("./%s/%s/token.json", account.LocalFolder, account.User))
		if utils.IsAfterCurrentTime(token.Expiry) {
			log.Println("Token expired getting new one")
			token = GetToken(account)
		}
	}

	return token

}

func GetToken(account structs.Account) (token *oauth2.Token) {
	googleOauthConfig := LoadConfiguration("secret.json")

	token, err := server.OauthServer(googleOauthConfig)
	if err != nil {
		log.Println(err)
	}
	WriteToken(token, account)
	return token
}
