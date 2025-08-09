package files

import (
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/api/drive/v3"
)

func DownloadFile(srv *drive.Service, id string) {
	resp, err := srv.Files.Get(id).Download()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	dFile, err := srv.Files.Get(id).Fields("name").Do()
	if err != nil {
		log.Fatal(err.Error())
	}

	file, err := os.Create("output/" + dFile.Name)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer file.Close()

	io.Copy(file, resp.Body)
}
