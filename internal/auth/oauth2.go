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

	"github.com/rogafe/go-email/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func LoadConfiguration(file string) (config *oauth2.Config) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err = google.ConfigFromJSON(b, "https://mail.google.com")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return config

}

func GoogleOauth() (token *oauth2.Token) {
	googleOauthConfig := LoadConfiguration("secret.json")

	oauthStateString := utils.GeneratePassword(10, 60)

	fmt.Printf("Please go to this address to get the token\n%s\n", googleOauthConfig.AuthCodeURL(oauthStateString))
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("copy the token here >>> ")
	userInput, _ := reader.ReadString('\n')
	token = requestToken(userInput, googleOauthConfig)

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
		panic(err)
	}

	return token
}
