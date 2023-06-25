package finder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/mod/modfile"
)

func New() *Model {
	return &Model{
		selected: -1,
	}
}

type goProject struct {
	path string
}

type Model struct {
	root     string
	choices  []*goProject
	selected int
	cursor   int
}

func (m *Model) SetCurrent(root string) {
	m.root = root

	m.choices = []*goProject{}
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.Name() == "go.mod" {
			m.choices = append(m.choices, &goProject{path})
		}
		return nil
	})
}
func (m *Model) GetSelected() string {
	if m.selected == -1 {
		return ""
	}
	return m.choices[m.selected].path
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", "l":
			m.selected = m.cursor
		}
	}

	return m, nil
}

func (m *Model) View() string {
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#363cb3")).Bold(true)

	s := ""
	for i, choice := range m.choices {
		cursor := " "
		path := choice.path
		if m.cursor == i {
			cursor = ">"
			path = selectedStyle.Render(path)
		}
		s += fmt.Sprintf("%s %s\n", cursor, path)
	}
	return s
}

func readModFile(path string) (*modfile.File, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f, err := modfile.Parse(path, content, nil)
	if err != nil {
		return nil, err
	}
	return f, nil
}
