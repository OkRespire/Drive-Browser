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

func (m *gModel) SaveCurrentState() {
	state := NavigationState{
		files:              make([]*drive.File, len(m.files)),
		pages:              make([][]*drive.File, len(m.pages)),
		pageCount:          m.pageCount,
		currentFolderId:    m.currentFolderId,
		nextPageToken:      m.nextPageToken,
		previousPageTokens: make([]string, len(m.previousPageTokens)),
		finalPage:          m.finalPage,
		cursor:             m.cursor,
	}

	copy(state.files, m.files)
	for i, page := range m.pages {
		state.pages[i] = make([]*drive.File, len(page))
		copy(state.pages[i], page)
	}

	copy(state.previousPageTokens, m.previousPageTokens)

	m.navigationStack = append(m.navigationStack, state)
}

func (m *gModel) OpenFolder(id string) error {
	r, err := m.srv.Files.List().PageSize(10).
		Q(fmt.Sprintf("'%s' in parents", id)).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		return err
	}

	m.SaveCurrentState()

	m.FindBreadCrumb(m.srv, m.files[m.cursor].Id)

	m.files = r.Files
	m.currentFolderId = id
	m.nextPageToken = r.NextPageToken
	m.cursor = 0
	m.pageCount = 1
	m.previousPageTokens = []string{}
	m.pages = [][]*drive.File{}

	return nil
}

func (m *gModel) LoadNextPage(pageToken string) error {
	call := m.srv.Files.List().PageSize(10).
		OrderBy("name").
		Fields("nextPageToken, files(id, name)")

	if pageToken != "" {
		call = call.PageToken(pageToken)
	}

	res, err := call.Do()
	if err != nil {
		return err
	}

	m.files = res.Files
	m.nextPageToken = res.NextPageToken
	m.cursor = 0
	m.pages = append(m.pages, res.Files)

	m.pageCount++
	return nil
}

func (m *gModel) LoadCachedPage(currPage int) error {

	index := currPage - 1

	if index < 0 {
		return nil
	}
	m.files = m.pages[index]
	m.pageCount--
	m.nextPageToken = m.previousPageTokens[index]
	m.cursor = 0

	return nil

}

func (m *gModel) RestorePreviousState() error {

	if len(m.navigationStack) == 0 {
		return fmt.Errorf("no previous state to restore")
	}

	lastIndex := len(m.navigationStack) - 1
	state := m.navigationStack[lastIndex]

	m.files = state.files
	m.pages = state.pages
	m.pageCount = state.pageCount
	m.currentFolderId = state.currentFolderId
	m.nextPageToken = state.nextPageToken
	m.previousPageTokens = state.previousPageTokens
	m.finalPage = state.finalPage
	m.cursor = state.cursor

	m.navigationStack = m.navigationStack[:lastIndex]

	// Remove last breadcrumb
	if len(m.breadcrumb) > 1 {
		m.breadcrumb = m.breadcrumb[:len(m.breadcrumb)-2]
	}

	return nil
}
