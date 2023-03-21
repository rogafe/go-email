package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func OauthServer() (token *oauth2.Token, err error) {
	// Load the Google OAuth secrets from the secret.json file
	b, err := os.ReadFile("secret.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read secret file: %v", err)
	}
	conf, err := google.ConfigFromJSON(b, "https://mail.google.com")
	if err != nil {
		return nil, fmt.Errorf("failed to parse secret file: %v", err)
	}

	// Start an HTTP server that listens on localhost:8080
	authCodeCh := make(chan string)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			authCodeCh <- code
			fmt.Fprint(w, "Authorization successful!")
		} else {
			url := conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
			http.Redirect(w, r, url, http.StatusFound)
		}
	})
	go http.ListenAndServe(":8080", nil)

	// Wait for the user to authorize the application and get the authorization code
	fmt.Println("Please visit the following URL to authorize the application:")
	url := conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println(url)
	openBrowser(url)
	code := <-authCodeCh

	// Exchange the authorization code for an access token
	token, err = conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}
	return token, nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch os := runtime.GOOS; os {
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	default:
		return fmt.Errorf("unsupported platform")
	}

	args = append(args, url)
	log.Printf("Opening url %s with default browser", url)
	return exec.Command(cmd, args...).Start()
}
