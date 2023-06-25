package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/mod/modfile"
)

const version = "day-2"

func main() {
	fmt.Printf("copper %s\n", version)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Current directory: %v\n", err)
		os.Exit(1)
	}

	choices, err := os.ReadDir(cwd)
	if err != nil {
		fmt.Printf("Reading directory: %v\n", err)
		os.Exit(1)
	}

	model := &copper{
		cwd:     cwd,
		choices: choices,
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

type copper struct {
	cwd     string
	choices []os.DirEntry
	cursor  int
	err     error
}

func (c *copper) Init() tea.Cmd { return nil }

func (c *copper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return c, tea.Quit

		case "up", "k":
			if c.cursor > 0 {
				c.cursor--
			}

		case "down", "j":
			if c.cursor < len(c.choices)-1 {
				c.cursor++
			}

		case "enter", "l":
			choice := c.choices[c.cursor]
			if !choice.Type().IsDir() {
				break
			}
			c.cwd = filepath.Join(c.cwd, c.choices[c.cursor].Name())
			c.choices, c.err = os.ReadDir(c.cwd)
			c.cursor = 0

		case "backspace", "h":
			c.cwd = filepath.Dir(c.cwd)
			c.choices, c.err = os.ReadDir(c.cwd)
			c.cursor = 0
		}
	}

	return c, nil
}

func (c *copper) View() string {
	if c.err != nil {
		return fmt.Sprintf("error: %v\n\nPress backspace to go back.\nPress q to quit.\n", c.err)
	}

	goStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#363cb3"))
	goBoldStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#363cb3")).Bold(true)
	hiddenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#999999"))
	dirStyle := lipgloss.NewStyle().Bold(true)

	var s string

	for i, choice := range c.choices {
		cursor := " "
		if c.cursor == i {
			cursor = ">"
		}

		name := choice.Name()

		if name == "go.mod" {
			gomod, err := readModFile(filepath.Join(c.cwd, name))
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

	s = fmt.Sprintf("%s\n\n", c.cwd) + s
	s += "\nPress q to quit.\n"

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
