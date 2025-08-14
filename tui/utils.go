package tui

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

func (m *gModel) FindBreadCrumb(srv *drive.Service, folderId string) error {
	f, err := srv.Files.Get(folderId).Fields("name").Do()
	if err != nil {
		return err
	}

	m.breadcrumb = append(m.breadcrumb, f.Name)

	return nil

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

	m.FindBreadCrumb(m.srv, m.files[m.cursor].Id)

	m.files = r.Files
	m.nextPageToken = r.NextPageToken
	m.cursor = 0
	m.pageCount = 1
	m.previousPageTokens = []string{}
	m.pages = [][]*drive.File{}

	return nil
}
