package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	oauth2 "golang.org/x/oauth2"
)

func tokenFromFile(file string) (*oauth2.Token, error) {
	// moste oppna filen
	f, err := os.Open(file)
	// om jag inte hittar filen, returnera fel.
	if err != nil {
		return nil, err
	}
	// måste stänga filen när jag func är klar
	defer f.Close()

	// skapa en ny tom token struct
	var token oauth2.Token
	err = json.NewDecoder(f).Decode(&token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func saveToken(path string, token *oauth2.Token) error {
	f, err := os.Create(path) // skapar // skriver ;ver filen p[ path

	if err != nil {
		// g[r det inte retunera fel
		return err
	}
	// st'nger filen n'r vi 'r klara
	defer f.Close()

	return json.NewEncoder(f).Encode(token) // skriver token som en JSON fil
}

// skapar en http klient fr'n en oauth2 config som ska användas för att göra anrop mot Google Drive API:et
func GetClient(config *oauth2.Config) *http.Client {
	// sparar Token
	tokenFil := "token.json"
	// försöker ladda/l'sa token fr'n fil
	token, err := tokenFromFile(tokenFil)

	if err != nil {
		// finns den inte så startar vi webbflödet för att få en ny token
		authenticateURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		// visar vilken url anv'dnaren ska bes;ka
		println("Gå till följande länk i din webbläsare och klistra in koden här:\n%v\n", authenticateURL)

		var code string
		if _, err := fmt.Scan(&code); err != nil {
			log.Fatalf("Kunde inte läsa in koden: %v", err)
		}
		// byter kod mot token
		tok, err := config.Exchange(context.TODO(), code)
		if err != nil {
			log.Fatalf("Kunde inte byta kod mot token: %v", err)
		}

		// sparar token för framtiden.
		saveToken(tokenFil, tok)
		token = tok
	}
	// skapar en http klient med token
	return config.Client(context.TODO(), token)
}
