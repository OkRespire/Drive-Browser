package tui

import (
	"google.golang.org/api/drive/v3"
)

func FindBreadCrumb(srv *drive.Service, folderId string) ([]string, error) {
	breadcrumb := []string{}
	currentId := folderId

	for {
		f, err := srv.Files.Get(currentId).Fields("id, name, parents").Do()
		if err != nil {
			return nil, err
		}

		breadcrumb = append([]string{f.Name}, breadcrumb...)

		if len(f.Parents) == 0 || f.Parents[0] == "root" {
			break
		}

		currentId = f.Parents[0]
	}

	return breadcrumb, nil

}
