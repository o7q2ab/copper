package picker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/mod/modfile"
)

func New() (*Model, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd(): %w", err)
	}

	choices, err := os.ReadDir(cwd)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir(): %w", err)
	}

	model := &Model{
		cwd:     cwd,
		choices: choices,
	}

	return model, nil
}

type Model struct {
	cwd     string
	choices []os.DirEntry
	cursor  int
	err     error
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
			choice := m.choices[m.cursor]
			if !choice.Type().IsDir() {
				break
			}
			m.cwd = filepath.Join(m.cwd, m.choices[m.cursor].Name())
			m.choices, m.err = os.ReadDir(m.cwd)
			m.cursor = 0

		case "backspace", "h":
			m.cwd = filepath.Dir(m.cwd)
			m.choices, m.err = os.ReadDir(m.cwd)
			m.cursor = 0
		}
	}

	return m, nil
}

func (m *Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("error: %v\n\nPress backspace to go back.\nPress q to quit.\n", m.err)
	}

	goStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#363cb3"))
	goBoldStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#363cb3")).Bold(true)
	hiddenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#999999"))
	dirStyle := lipgloss.NewStyle().Bold(true)

	var s string

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		name := choice.Name()

		if name == "go.mod" {
			gomod, err := readModFile(filepath.Join(m.cwd, name))
			if err != nil {
				s = fmt.Sprintf("Module: (error) %s\n\n", err)
			}
			directDeps, indirectDeps := 0, 0
			for _, d := range gomod.Require {
				if d.Indirect {
					indirectDeps++
				} else {
					directDeps++
				}
			}
			s = fmt.Sprintf("Module: %s\n\tdirect dependencies: %d\n\tindirect dependencies: %d\n\n",
				gomod.Module.Mod.Path, directDeps, indirectDeps) + s
		}

		switch {
		case name == "go.mod" || name == "go.sum" || name == "go.work" || name == "go.work.sum":
			name = goBoldStyle.Render(name)

		case strings.HasSuffix(name, ".go"):
			name = goStyle.Render(name)

		case strings.HasPrefix(name, "."):
			name = hiddenStyle.Render(name)
		}

		if choice.Type().IsDir() {
			name = dirStyle.Render(name)

			s += fmt.Sprintf("%s %s\n", cursor, name)
		} else {
			info, err := choice.Info()
			if err != nil {
				s += fmt.Sprintf("%s %s (%v)\n", cursor, name, err)
			} else {
				s += fmt.Sprintf("%s %s (%d bytes)\n", cursor, name, info.Size())
			}
		}
	}

	s = fmt.Sprintf("%s\n\n", m.cwd) + s

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
