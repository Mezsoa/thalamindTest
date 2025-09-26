package main

import "os"








// DownloadGoogleDriveFolder laddar ner alla filer och undermappar i en Google Drive-mapp rekursivt.
// anv'nder sem f;r att begr'nnsa antalet samtidiga nedladdningar s[ att inte google API:et blockeras. och kraschar.
func DownloadGoogleDriveFolder( srv *drive.Service, folderID, localPath string, wg. *sync.WaitGroup, sem chan struct{}) error {
	
	// funkar som signa; f;r att denna go routine 'r klar och att funktionen har avslutats
	defer wg.Done()

	if err := os.MkdirAll()
}
