package files

import (
	"log"

	"google.golang.org/api/drive/v3"
)

func ListFiles(srv *drive.Service) []*drive.File {
	r, err := srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	return r.Files
}
