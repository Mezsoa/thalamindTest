package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func main() {

	// läser in credentials.json filen
	data, err := os.ReadFile("credentialsThalamind.json")
	if err != nil {
		log.Fatalf("Kunde inte läsa credentials.json: %v", err)
	}
	// skapar en OAuth2 config från credentials filen
	config, err := google.ConfigFromJSON(data, drive.DriveReadonlyScope)
	if err != nil {
		log.Fatalf("Kunde inte parsa credentials.json till configen: %v", err)
	}

	// H'mta/ får en authad http klient
	client := GetClient(config)

	// skapar en Drive service fr'n klienten
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Kunde inte skapa Drive service: %v", err)
	}

	// Go
	home, _ := os.UserHomeDir()

	// ID fför mappen på google drive som jag vill ladda ner
	rootFolderID := "1HjVNaFN2G0jIsE8T12CAE94aqE9czNmN"
	localPath := filepath.Join(home, "Desktop", "GoogleDriveFolder")

	// skapar en waitgroup för att vänta på att alla go routines ska bli klara
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // sem kanal f;r att begränsa antalet samtidiga nedladdningar s'tter den till 5

	// startar nedladdningen av mappen
	wg.Add(1)
	go DownloadGoogleDriveFolder(srv, rootFolderID, localPath, &wg, sem)

	// väntar på att alla nedladdningar ska bli klara
	wg.Wait()
	log.Println("Nedladdning klar!")
}
