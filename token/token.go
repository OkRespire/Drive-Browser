package token

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"

	"golang.org/x/oauth2"
)

func getTokenFromLocalServer(config *oauth2.Config) *oauth2.Token {
	mux := http.NewServeMux()
	listener, err := net.Listen("tcp", "localhost:8080")
	fmt.Println("Listening...")
	if err != nil {
		log.Fatalf("Unable to start local server: %v", err)
	}

	codeCh := make(chan string)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Authorisation in progress")
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "No code in request", http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, "Authorization complete. You can close this window.")
		codeCh <- code
	})
	server := &http.Server{
		Handler: mux,
	}
	go func() {
		err := http.Serve(listener, nil)
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	openURL(authURL)
	var code string
	if runtime.GOOS == "windows" {
		fmt.Println("A new browser tab will be opened.\nPaste the code found in the url here:")
		fmt.Println("The part of the URL that  you would need to focus on is code=YOUR_CODE. Ignore the & at the end")
		fmt.Print("Code: ")
		fmt.Scanln(&code)
	} else {
		code = <-codeCh
	}
	server.Shutdown(context.Background())

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}
	return token
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "./token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromLocalServer(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}
