package tui

import (
	"fmt"

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

func (m gModel) MimeTypeCheck(id string) (string, error) {

	dFile, err := m.srv.Files.Get(id).Fields("mimeType").Do()
	if err != nil {
		return "", err
	}

	return dFile.MimeType, nil

}

func (m *gModel) OpenFolder(id string) error {
	r, err := m.srv.Files.List().PageSize(10).
		Q(fmt.Sprintf("'%s' in parents", id)).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		return err
	}

	m.files = r.Files
	m.nextPageToken = r.NextPageToken
	m.cursor = 0
	m.pageCount = 1
	m.previousPageTokens = []string{}
	m.pages = [][]*drive.File{}

	return nil
}
