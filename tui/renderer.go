package tui

import (
	"context"
	"fmt"
	"log"

	"drivebrowser/files"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"google.golang.org/api/drive/v3"
)

type gModel struct {
	breadcrumb []string
	files      []*drive.File
	cursor     int
	user       *drive.User
	srv        *drive.Service
}

func InitialModel(ctx context.Context, srv *drive.Service, folderId string) gModel {
	file_list := files.ListFiles(srv)
	breadcrumbVal, err := FindBreadCrumb(srv, folderId)
	if err != nil {
		log.Fatal(err.Error())
	}

	user_name, err := srv.About.Get().Fields("user(displayName, emailAddress)").Context(ctx).Do()
	if err != nil {
		log.Fatal(err.Error())
	}

	return gModel{
		breadcrumb: breadcrumbVal,
		files:      file_list,
		cursor:     0,
		user:       user_name.User,
		srv:        srv,
	}
}

func (m gModel) Init() tea.Cmd {
	return nil
}

func (m gModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

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
		case "enter":
			files.DownloadFile(m.srv, m.files[m.cursor].Id)
		}
	}

	return m, nil
}

func (m gModel) View() string {

	breadcrumbStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Border(lipgloss.NormalBorder(), true).
		Padding(0, 1)

	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9AA0A6")).
		Border(lipgloss.NormalBorder(), true).
		Padding(0, 1).Bold(true)
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

		file_string += fmt.Sprintf("\n%s %s", cursor, f.Name)
	}

	// Layout
	breadcrumbBar := breadcrumbStyle.Render(breadcrumb_string)
	content := contentStyle.Render(file_string)

	return lipgloss.JoinVertical(lipgloss.Center,
		breadcrumbBar,
		content,
	)
}
