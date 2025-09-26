package main

import (
	"context"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func main() {

	// läser in credentials.json filen
	data, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Kunde inte läsa credentials.json: %v", err)
	}
	// Oauth2 config frn credentials filen
	config, err := google.ConfigFromJSON(data, drive.DriveReadonlyScope)

	if err != nil {
		log.Fatalf("Kunde inte parsa credentials.json till configen: %v", err)
	}

	// H'mta en authad http klient
	client := getClient(config)

	// skapar en Drive service fr'n klienten
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))

	if err != nil {
		log.Fatalf("Kunde inte skapa Drive service: %v", err)
	}

}
