package tui

import (
	"log"

	"drivebrowser/files"

	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/api/drive/v3"
)

type gModel struct {
	breadcrumb []string
	files      []*drive.File
	cursor     int
}

func InitialModel(srv *drive.Service, folderId string) gModel {
	file_list := files.ListFiles(srv)
	breadcrumbVal, err := FindBreadCrumb(srv, folderId)
	if err != nil {
		log.Fatal(err.Error())
	}

	return gModel{
		breadcrumb: breadcrumbVal,
		files:      file_list,
		cursor:     0,
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
			if m.cursor > 0 {
				m.cursor--
			}

		}
	}

	return m, nil
}

func (m gModel) View() string {
	s := ""

	for _, v := range m.breadcrumb {
		s += v
	}

	for _, f := range m.files {
		s += "\n" + f.Name
	}

	return s
}
