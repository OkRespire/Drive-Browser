package tui

import (
	"context"
	"errors"
	"fmt"
	"log"

	"drivebrowser/files"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"google.golang.org/api/drive/v3"
)

type gModel struct {
	breadcrumb         []string
	files              []*drive.File
	pages              [][]*drive.File
	cursor             int
	user               *drive.User
	srv                *drive.Service
	pageCount          int
	nextPageToken      string
	previousPageTokens []string
	finalPage          bool
	width              int
	height             int
}

func InitialModel(ctx context.Context, srv *drive.Service, folderId string) gModel {
	file_list, nextToken := files.ListFiles(srv)
	breadcrumbVal, err := FindBreadCrumb(srv, folderId)
	if err != nil {
		log.Fatal(err.Error())
	}

	user_name, err := srv.About.Get().Fields("user(displayName, emailAddress)").Context(ctx).Do()
	if err != nil {
		log.Fatal(err.Error())
	}

	return gModel{
		breadcrumb:         breadcrumbVal,
		files:              file_list,
		pages:              [][]*drive.File{file_list},
		cursor:             0,
		user:               user_name.User,
		srv:                srv,
		pageCount:          1,
		nextPageToken:      nextToken,
		previousPageTokens: []string{},
		finalPage:          false,
		width:              0,
		height:             0,
	}
}

func (m *gModel) loadNextPage(pageToken string) error {
	call := m.srv.Files.List().PageSize(10).
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

func (m gModel) Init() tea.Cmd {
	return nil
}

func (m gModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.files) - 1
			}
		case "down", "j":
			m.cursor++
			if m.cursor > len(m.files)-1 {
				m.cursor = 0
			}
		case "right", "l":
			if m.nextPageToken != "" {
				m.previousPageTokens = append(m.previousPageTokens, m.nextPageToken)
				err := m.loadNextPage(m.nextPageToken)
				if err != nil {
					log.Fatal("Error loading next page:", err)
				}
			} else {
				m.finalPage = true
			}

		case "left", "h":
			m.finalPage = false
			if len(m.previousPageTokens) > 0 {
				err := m.loadCachedPage(m.pageCount - 1)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		case "enter":
			mimeType, err := m.MimeTypeCheck(m.files[m.cursor].Id)
			if err != nil {
				log.Fatal(err.Error())
			}
			if mimeType == "application/vnd.google-apps.folder" {
				fmt.Println(m.files[m.cursor].Name)
				m.OpenFolder(m.files[m.cursor].Id)

			} else {
				files.DownloadFile(m.srv, m.files[m.cursor].Id)
			}
		}
	}

	return m, nil
}

func (m *gModel) loadCachedPage(currPage int) error {

	index := currPage - 1

	if index < 0 {
		return errors.New("index out of bounds")
	}
	m.files = m.pages[index]
	m.pageCount--
	m.nextPageToken = m.previousPageTokens[index]
	m.cursor = 0

	return nil

}

func (m gModel) View() string {
	maxNameLen := 0
	for _, f := range m.files {
		if len(f.Name) > maxNameLen {
			maxNameLen = len(f.Name)
		}
	}

	breadcrumbStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Border(lipgloss.NormalBorder(), true).
		Padding(0, 1)

	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9AA0A6")).
		Padding(0, 0, 1).Bold(true)

	pageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9AA0A6")).
		Padding(0, 0, 1).Bold(true)

	breadcrumb_string := fmt.Sprintf("%s (%s)\n", m.user.DisplayName, m.user.EmailAddress)

	for _, v := range m.breadcrumb {
		breadcrumb_string += v
	}

	file_string := ""

	for i, f := range m.files {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		file_string += fmt.Sprintf("\n%s %-*s", cursor, maxNameLen, f.Name)
	}

	var page_string string
	if m.finalPage {
		page_string = fmt.Sprintf("< Page: %d", m.pageCount)
	} else if m.pageCount == 1 {
		page_string = fmt.Sprintf("Page: %d >", m.pageCount)
	} else {
		page_string = fmt.Sprintf("< Page: %d >", m.pageCount)
	}

	// Layout
	breadcrumbBar := lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Center,
		breadcrumbStyle.Render(breadcrumb_string),
	)
	content := contentStyle.Render(file_string)
	page := lipgloss.PlaceHorizontal(
		lipgloss.Width(breadcrumbBar),
		lipgloss.Center,
		pageStyle.Render(page_string),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		breadcrumbBar,
		content,
		page,
	)
}
