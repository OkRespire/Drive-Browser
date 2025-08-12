package files

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/drive/v3"
)

var mimeTypes = map[string]string{
	"application/vnd.google-apps.document":     "application/pdf",
	"application/vnd.google-apps.spreadsheet":  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"application/vnd.google-apps.presentation": "application/pdf",
	"application/vnd.google-apps.drawing":      "image/png",
}

func DownloadFile(srv *drive.Service, id string) {
	dFile, err := srv.Files.Get(id).Fields("name, mimeType").Do()
	if err != nil {
		log.Fatal(err.Error())
	}
	var resp *http.Response
	if k, v := mimeTypes[dFile.MimeType]; v {
		resp, err = srv.Files.Export(id, k).Download()

		if err != nil {
			fmt.Println(err.Error())
		}

	} else {
		resp, err = srv.Files.Get(id).Download()

		if err != nil {
			fmt.Println(err.Error())
		}

	}
	defer resp.Body.Close()
	if err := os.MkdirAll("output", 0755); err != nil {
		log.Fatal(err.Error())
	}

	file, err := os.Create("output/" + dFile.Name)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer file.Close()

	io.Copy(file, resp.Body)
}
