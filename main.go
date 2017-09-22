package main

import (
	"fmt"
	"github.com/sequoiia/gtarma-dl/client"
	"net/http"
	"os"
)

var savePath string = ""

func main() {
	fmt.Println("Arma 3 mission downloader - @sequoiia")

	savePath = os.Getenv("LOCALAPPDATA") + "\\Arma 3\\MPMissionsCache\\"

	c := client.DlClient{HttpClient: http.DefaultClient}
	files := c.GetFiles()
	c.DownloadFiles(files, savePath)

	fmt.Println("Download finished.")
}
