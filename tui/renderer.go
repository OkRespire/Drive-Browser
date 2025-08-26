package tui

import (
	"context"
	"fmt"
	"log"

	"drivebrowser/files"
	"drivebrowser/utils"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"google.golang.org/api/drive/v3"
)

func InitialModel(ctx context.Context, srv *drive.Service, folderId string) gModel {
	file_list, nextToken := files.ListFiles(srv)

	user_name, err := srv.About.Get().Fields("user(displayName, emailAddress)").Context(ctx).Do()
	if err != nil {
		log.Fatal(err.Error())
	}

	return gModel{
		breadcrumb:         []string{"My Drive"},
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
		currentFolderId:    folderId,
		navigationStack:    []NavigationState{},
		isSearching:        false,
		searchQuery:        "",
	}
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
		if m.isSearching {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				if m.searchQuery != "" {
					err := m.Search()
					if err != nil {
						log.Fatal("Search error:", err)
					}
				}
				m.isSearching = false
				m.searchQuery = ""
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				}

			case "esc":
				m.isSearching = false

			default:
				if len(msg.String()) == 1 {
					m.searchQuery += msg.String()
				}

			}
			return m, nil

		}

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
				err := m.LoadNextPage(m.nextPageToken)
				if err != nil {
					log.Fatal("Error loading next page:", err)
				}
			} else {
				m.finalPage = true
			}

		case "left", "h":
			m.finalPage = false
			if len(m.previousPageTokens) > 0 {
				err := m.LoadCachedPage(m.pageCount - 1)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		case "enter":
			mimeType := m.files[m.cursor].MimeType
			if mimeType == "application/vnd.google-apps.folder" {
				m.OpenFolder(m.files[m.cursor].Id)

			} else {
				files.DownloadFile(m.srv, m.files[m.cursor].Id)
			}
		case "backspace":
			if err := m.RestorePreviousState(); err != nil {
				log.Fatal(err.Error())
			}
		case "/":
			m.isSearching = true
			m.searchQuery = ""

		}
	}

	return m, nil
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

	for i, v := range m.breadcrumb {
		if i > 0 {
			breadcrumb_string += " > "
		}
		breadcrumb_string += v
	}

	file_string := ""

	for i, f := range m.files {
		icon := utils.GetFileIcon(f.Name, f.MimeType)

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		file_string += fmt.Sprintf("\n%s %s %-*s", cursor, icon, maxNameLen, f.Name)
	}

	var page_string string
	if m.finalPage {
		page_string = fmt.Sprintf("< Page: %d", m.pageCount)
	} else if m.pageCount == 1 {
		page_string = fmt.Sprintf("Page: %d >", m.pageCount)
	} else {
		page_string = fmt.Sprintf("< Page: %d >", m.pageCount)
	}
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

	if m.isSearching {
		// Show search input at bottom
		searchInput := fmt.Sprintf("Search: %s_", m.searchQuery)
		return lipgloss.JoinVertical(lipgloss.Left,
			breadcrumbBar,
			content,
			page,
			lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render(searchInput),
		)
	}
	// Layout

	return lipgloss.JoinVertical(lipgloss.Left,
		breadcrumbBar,
		content,
		page,
	)
}
