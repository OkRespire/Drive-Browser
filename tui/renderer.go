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
		searchModel:        nil,
		isTyping:           false,
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
					m.isTyping = false
				}
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				}

			case "esc":
				m.isSearching = false
				m.searchQuery = ""
			case "/":
				m.searchQuery = ""
				m.isTyping = true

			default:
				if len(msg.String()) == 1 {
					m.searchQuery += msg.String()
				}

			}

		}
		var currentFiles []*drive.File
		var currentCursor *int

		if m.searchModel != nil {
			currentFiles = m.searchModel.files
			currentCursor = &m.searchModel.cursor
		} else {
			currentFiles = m.files
			currentCursor = &m.cursor
		}
		if len(currentFiles) == 0 {
			return m, nil
		}

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			*currentCursor--
			if *currentCursor < 0 {
				*currentCursor = len(currentFiles) - 1
			}
		case "down", "j":
			*currentCursor++
			if *currentCursor > len(currentFiles)-1 {
				*currentCursor = 0
			}
		case "right", "l":
			if m.isSearching && m.searchModel != nil {
				if m.searchModel.nextPageToken != "" {
					err := m.LoadNextSearchPage()
					if err != nil {
						log.Fatal("Error loading next search page:", err)
					}
				} else {
					m.searchModel.finalPage = true
				}
			} else {
				if m.nextPageToken != "" {
					err := m.LoadNextPage()
					if err != nil {
						log.Fatal("Error loading next page:", err)
					}
				} else {
					m.finalPage = true
				}
			}

		case "left", "h":
			if m.isSearching && m.searchModel != nil {
				m.searchModel.finalPage = false
				if len(m.searchModel.previousPageTokens) > 0 {
					err := m.LoadSearchCachedPage(m.searchModel.pageCount - 1)
					if err != nil {
						log.Fatal(err.Error())
					}
				}
			} else {
				m.finalPage = false
				if len(m.previousPageTokens) > 0 {
					err := m.LoadCachedPage(m.pageCount - 1)
					if err != nil {
						log.Fatal(err.Error())
					}
				}
			}
		case "enter":
			mimeType := currentFiles[*currentCursor].MimeType
			if mimeType == "application/vnd.google-apps.folder" {
				m.OpenFolder(currentFiles[*currentCursor].Id)

			} else {
				files.DownloadFile(m.srv, currentFiles[*currentCursor].Id)
			}
		case "backspace":
			if err := m.RestorePreviousState(); err != nil {
				log.Fatal(err.Error())
			}
		case "/":
			m.isSearching = true
			m.searchQuery = ""
			m.searchModel = nil
			m.isTyping = true

		}
	}

	return m, nil
}

func (m gModel) View() string {
	var files []*drive.File
	var cursorNum int
	var pageCount int
	var finalPage bool
	var breadcrumb []string

	breadcrumb = m.breadcrumb
	if m.isSearching && m.searchModel != nil {
		files = m.searchModel.files
		cursorNum = m.searchModel.cursor
		pageCount = m.searchModel.pageCount
		finalPage = m.searchModel.finalPage
	} else {
		files = m.files
		cursorNum = m.cursor
		pageCount = m.pageCount
		finalPage = m.finalPage
	}
	maxNameLen := 0
	for _, f := range files {
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

	for i, v := range breadcrumb {
		if i > 0 {
			breadcrumb_string += " > "
		}
		breadcrumb_string += v
	}

	file_string := ""

	for i, f := range files {
		icon := utils.GetFileIcon(f.Name, f.MimeType)

		cursor := " "
		if cursorNum == i {
			cursor = ">"
		}

		file_string += fmt.Sprintf("\n%s %s %-*s", cursor, icon, maxNameLen, f.Name)
	}

	var page_string string
	if finalPage {
		page_string = fmt.Sprintf("< Page: %d", pageCount)
	} else if pageCount == 1 {
		page_string = fmt.Sprintf("Page: %d >", pageCount)
	} else {
		page_string = fmt.Sprintf("< Page: %d >", pageCount)
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

	if m.isTyping {
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
