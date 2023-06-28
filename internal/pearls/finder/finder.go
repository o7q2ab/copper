package finder

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func New() *Model {
	return &Model{}
}

type goProject struct {
	path string
}

type Model struct {
	root    string
	choices []*goProject
	cursor  int
}

func (m *Model) SetCurrent(root string) {
	m.root = root
	m.cursor = 0
	m.choices = []*goProject{}
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		switch filepath.Base(path) {
		case "go.mod":
			m.choices = append(m.choices, &goProject{path})
		case ".git":
			return filepath.SkipDir
		case "vendor":
			return filepath.SkipDir
		}
		return nil
	})
}
func (m *Model) GetSelected() string {
	return m.choices[m.cursor].path
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

		case "backspace", "h":
			m.SetCurrent(filepath.Dir(m.root))
		}
	}

	return m, nil
}

func (m *Model) View() string {
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#363cb3")).Bold(true)

	s := fmt.Sprintf("%s\n", m.root)
	for i, choice := range m.choices {
		cursor := " "
		path := strings.TrimPrefix(choice.path, m.root)
		if m.cursor == i {
			cursor = ">"
			path = selectedStyle.Render(path)
		}
		s += fmt.Sprintf("%s %s\n", cursor, path)
	}
	return s
}
