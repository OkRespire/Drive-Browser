package tui

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

type searchModel struct {
	files              []*drive.File
	pages              [][]*drive.File
	cursor             int
	pageCount          int
	nextPageToken      string
	previousPageTokens []string
	finalPage          bool
}

func (m *gModel) LoadNextSearchPage() error {
	q := fmt.Sprintf("name contains '%s'", m.searchQuery)

	call := m.srv.Files.List().
		PageSize(10).
		OrderBy("name").
		Q(q).
		Fields("nextPageToken, files(id, name, mimeType)")

	if m.searchModel.nextPageToken != "" {
		call = call.PageToken(m.searchModel.nextPageToken)
	}

	res, err := call.Do()
	if err != nil {
		return err
	}

	m.searchModel.previousPageTokens = append(m.searchModel.previousPageTokens, m.searchModel.nextPageToken)
	m.searchModel.files = res.Files
	m.searchModel.nextPageToken = res.NextPageToken
	m.searchModel.cursor = 0
	m.searchModel.pages = append(m.searchModel.pages, res.Files)
	m.searchModel.pageCount++

	return nil
}

func (m *gModel) LoadSearchCachedPage(currPage int) error {
	index := currPage - 1

	if index < 0 || index >= len(m.searchModel.pages) {
		return nil
	}
	m.searchModel.files = m.searchModel.pages[index]
	m.searchModel.pageCount--
	m.searchModel.nextPageToken = m.searchModel.previousPageTokens[index]
	m.searchModel.cursor = 0

	return nil

}

func (m *gModel) SaveSearchModel(r *drive.FileList) {
	m.searchModel = &searchModel{
		files:              r.Files,
		pages:              [][]*drive.File{r.Files},
		cursor:             0,
		pageCount:          1,
		nextPageToken:      r.NextPageToken,
		previousPageTokens: []string{},
		finalPage:          r.NextPageToken == "",
	}

}

func (m *gModel) Search() error {

	r, err := m.srv.Files.List().PageSize(10).
		OrderBy("name").
		Q(fmt.Sprintf("name contains '%s'", m.searchQuery)).
		Fields("nextPageToken, files(id, name, mimeType)").Do()

	if err != nil {
		m.RestorePreviousState()
		return err
	}
	m.SaveSearchModel(r)

	m.isSearching = true

	return nil
}
