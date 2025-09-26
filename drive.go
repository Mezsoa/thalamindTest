package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"google.golang.org/api/drive/v3"
)

// DownloadGoogleDriveFolder laddar ner alla filer och undermappar i en Google Drive-mapp rekursivt.
// anv'nder sem f;r att begr'nnsa antalet samtidiga nedladdningar s[ att inte google API:et blockeras. och kraschar.
func DownloadGoogleDriveFolder(srv *drive.Service, folderID, localPath string, wg *sync.WaitGroup, sem chan struct{}) error {

	// funkar som signa; f;r att denna go routine 'r klar och att funktionen har avslutats
	defer wg.Done()

	if err := os.MkdirAll(localPath, 0755); err != nil {
		return fmt.Errorf("kunde inte skapa mapp %s: %v", localPath, err)
	}

	// listar alla mappar och filer i den drive mapp som jag pekar po
	query := fmt.Sprintf("'%s' in parents and trashed = false", folderID)
	fileList, err := srv.Files.List().Q(query).Fields("files(id, name, mimeType)").Do()

	if err != nil {
		return fmt.Errorf("kunde inte lista filer: %v", err)
	}

	for _, i := range fileList.Files {
		if i.MimeType == "application/vnd.google-apps.folder" {
			// om det 'r en mapp, skapa en ny lokal mapp och kalla p[ funktionen rekursivt
			subFolderPath := filepath.Join(localPath, i.Name)
			wg.Add(1)
			go func() {
				err := DownloadGoogleDriveFolder(srv, i.Id, subFolderPath, wg, sem)
				if err != nil {

				}
			}()
		} else {
			// om det 'r en fil, ladda ner den
			wg.Add(1)
			go func(file *drive.File) {
				defer wg.Done()
				sem <- struct{}{}        // skicka in i sem kanalen f;r att begr'nnsa samtidiga nedladdningar
				defer func() { <-sem }() // ta bort fr'n sem kanalen n'r nedladdningen 'r klar
				fmt.Println("nerladdade filer: ", file.Name)
				if err := downloadFile(srv, file.Id, filepath.Join(localPath, file.Name)); err != nil {
					fmt.Printf("Fel vid nedladdning av %s: %v\n", file.Name, err)
				}
			}(i)

		}
	}
	return nil
}

// downloadFile laddar ner en enskild fil fr'n Google Drive och sparar den lokalt
func downloadFile(srv *drive.Service, fileID, localPath string) error {
	// anropar Google Drive API:et f;r att kunna h'mta filen
	response, err := srv.Files.Get(fileID).Download()
	// blir det fel, returnera fel
	if err != nil {
		return err
	}
	// st'nger response n'r vi 'r klara
	defer response.Body.Close()

	// skapar en lokal fil f;r att spara inneh'lllet
	outFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	// st'nger filen n'r vi 'r klara
	defer outFile.Close()

	// kopierar response bodyns inneh[ll till den lokala filen outFile som jag skapade ovan
	_, err = io.Copy(outFile, response.Body)
	return err
}
