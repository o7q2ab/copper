package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

const version = "day-0"

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

	s := fmt.Sprintf("%s\n\n", c.cwd)

	for i, choice := range c.choices {
		cursor := " "
		if c.cursor == i {
			cursor = ">"
		}

		if choice.Type().IsDir() {
			s += fmt.Sprintf("%s %s\n", cursor, choice.Name())
		} else {
			info, err := choice.Info()
			if err != nil {
				s += fmt.Sprintf("%s %s (%v)\n", cursor, choice.Name(), err)
			} else {
				s += fmt.Sprintf("%s %s (%d bytes)\n", cursor, choice.Name(), info.Size())
			}
		}

	}

	s += "\nPress q to quit.\n"

	return s
}
