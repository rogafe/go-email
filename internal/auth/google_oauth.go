package auth

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

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
	tokenJson, err := json.MarshalIndent(token, "", "\t")
	if err != nil {
		log.Println(err)
	}

	utils.CreateFolder(fmt.Sprintf("./%s/%s", account.LocalFolder, account.User))
	err = ioutil.WriteFile(fmt.Sprintf("./%s/%s/token.json", account.LocalFolder, account.User), tokenJson, 0777)
	if err != nil {
		log.Println(err)
	}
}

func GoogleOauth(account structs.Account) (token *oauth2.Token) {

	exist, err := utils.FileExists(fmt.Sprintf("./%s/%s/token.json", account.LocalFolder, account.User))
	if err != nil {
		log.Println(err)
	}
	if !exist {
		googleOauthConfig := LoadConfiguration("secret.json")

		oauthStateString := utils.GeneratePassword(10, 60)

		fmt.Printf("Please go to this address to get the token\n%s\n", googleOauthConfig.AuthCodeURL(oauthStateString))
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("copy the token here >>> ")
		userInput, _ := reader.ReadString('\n')
		token = requestToken(userInput, googleOauthConfig)
		WriteToken(token, account)
	} else {
		token = LoadToken(fmt.Sprintf("./%s/%s/token.json", account.LocalFolder, account.User))
	}

	return token

}

func requestToken(userInput string, googleOauthConfig *oauth2.Config) (token *oauth2.Token) {
	u := "https://accounts.google.com/o/oauth2/token"

	data := url.Values{
		"code":          {userInput},
		"client_id":     {googleOauthConfig.ClientID},
		"client_secret": {googleOauthConfig.ClientSecret},
		"redirect_uri":  {googleOauthConfig.RedirectURL},
		"grant_type":    {"authorization_code"},
	}

	resp, err := http.PostForm(u, data)

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Panic(err)
	}
	log.Println(token)
	return token
}
